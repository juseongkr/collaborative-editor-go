package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/juseongkr/collaborative-editor-go/ot"
	"github.com/juseongkr/collaborative-editor-go/ot/server"
	"strconv"
)

type StateStore struct {
	name     string
	opStream chan server.OpMessage
}

func NewStateStore(name string) *StateStore {
	s := &StateStore{
		name:     name,
		opStream: make(chan server.OpMessage, 128),
	}

	go s.streamOperations()

	return s
}

func (s *StateStore) streamOperations() {
	streams := []string(s.name)
	lastID := "$"
	for {
		streams, err := rdb.XRead(&redis.XReadArgs{
			Streams: append(streams, lastID),
		}).Result()

		if err != nil {
			panic(err)
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var op ot.CompositeOp
				if err := json.Unmarshal([]byte(message.Values["op"].(string)), &op); err != nil {
					panic(err)
				}

				s.opStream <- server.OpMessage{
					Revision: idToRev(message.ID),
					Op:       op,
					AuthorID: message.Values["author"].(string),
				}
				lastID = message.ID
			}
		}
	}
}

func idToRev(id string) int {
	i, err := strconv.Atoi(id[2:])
	if err != nil {
		panic(err)
	}

	return i
}

func revToID(rev int) string {
	return "0-" + strconv.Itoa(rev)
}

func (s *StateStore) Current() (document string, revision int, err error) {
	res, err := rdb.XRevRangeN(s.name, "+", "-", 1).Result()
	if err != nil || len(res) == 0 {
		return
	}

	document = res[0].Values["doc"].(string)
	revision = idToRev(res[0].ID)

	return
}

func (s *StateStore) applyClientXAddArgs(opMsg server.OpMessage) (*redis.XAddArgs, error) {
	redisOps, err := rdb.XRange(s.name, revToID(opMsg.Revision), "+").Result()
	if err != nil {
		return nil, err
	}

	l := len(redisOps)
	if l > 0 {
		l--
	}

	ops := make([]ot.CompositeOp, l)
	if len(redisOps) > 0 {
		for i, redisOp := range redisOps[1:] {
			var op ot.CompositeOp
			if err := json.Unmarshal(
				[]byte(redisOp.Values["op"].(string)), &op); err != nil {
				return nil, err
			}

			ops[i] = op
		}
	}

	var currentDoc string
	var currentRev int
	if len(redisOps) > 0 {
		lastOp := redisOps[len(redisOps)-1]
		currentDoc = lastOp.Values["doc"].(string)
		currentRev = idToRev(lastOp.ID)
	}

	res, err := server.ApplyClientOp(server.ApplyClientOpInput{
		CurrentDocument: currentDoc,
		CurrentRevision: currentRev,
		Op:              opMsg.Op,
		TransformOps:    ops,
	})

	if err != nil {
		return nil, err
	}

	jsonOp, err := json.Marshal(res.Op)
	if err != nil {
		return nil, err
	}

	return &redis.XAddArgs{
		Stream: s.name,
		ID:     revToID(res.Revision),
		Values: map[string]interface{}{
			"op":     string(jsonOp),
			"doc":    res.Document,
			"author": opMsg.AuthorID,
		},
	}, nil

}

func (s *StateStore) ApplyClient(opMsg server.OpMessage) error {
	for {
		args, err := s.applyClientXAddArgs(opMsg)
		if err != nil {
			return err
		}

		if err := rdb.XAdd(args).Err(); err == nil {
			return nil
		}
	}
}

func (s *StateStore) OperationStream() <-chan server.OpMessage {
	return s.opStream
}
