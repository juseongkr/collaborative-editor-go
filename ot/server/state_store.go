package server

import (
	"github.com/juseongkr/collaborative-editor-go/ot"
	"sync"
)

type StateStore interface {
	Current() (document string, revision int, err error)
	ApplyClient(opMsg OpMessage) error
	OperationStream() <-chan OpMessage
}

type MemoryStateStore struct {
	mux      sync.RWMutex
	document string
	ops      []ot.CompositeOp
	opStream chan OpMessage
}

func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{
		opStream: make(chan OpMessage, 128),
	}
}

func (m *MemoryStateStore) Current() (document string, revision int, err error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	document = m.document
	revision = len(m.ops)

	return
}

func (m *MemoryStateStore) OperationStream() <-chan OpMessage {
	return m.opStream
}

func (m *MemoryStateStore) ApplyClient(opMsg OpMessage) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if opMsg.Revision < 0 || (opMsg.Revision > len(m.ops)) {
		return ErrUnknownRevision
	}

	res, err := ApplyClientOp(ApplyClientOpInput{
		CurrentDocument: m.document,
		CurrentRevision: len(m.ops),
		Op:              opMsg.Op,
		TransformOps:    m.ops[opMsg.Revision:],
	})

	if err != nil {
		return err
	}

	m.document = res.Document
	m.ops = append(m.ops, res.Op)

	m.opStream <- OpMessage{
		AuthorId: opMsg.AuthorID,
		Op:       res.Op,
		Revision: res.Revision,
	}

	return nil
}
