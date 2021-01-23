package server

import (
	"github.com/juseongkr/collaborative-editor-go/ot"
	"errors"
	"log"
)

var ErrUnknownRevision = errors.New("unknown revision")
var ErrInvalidOperation = errors.New("invalid operation")

type ApplyClientOpInput struct {
	CurrentDocument string
	CurrentRevision int
	Op              ot.CompositeOp
	TransformOps    []ot.CompositeOp
}

type ApplyClientOpOutput struct {
	Document string
	Op       ot.CompositeOp
	Revision int
}

func ApplyClientOp(i ApplyClientOpInput) (o ApplyClientOpOutput, err error) {
	log.Printf("input: %+v\n", i)
	o.Op = i.Op
	for _, transformOp := range i.TransformOps {
		if o.Op.InputLength() != transformOp.InputLength() {
			err = ErrInvalidOperation
			return
		}
		o.Op, _ = o.Op.Transform(transformOp)
	}

	newDoc, err := ot.ApplyString(o.Op, i.CurrentDocument)
	if err != nil {
		return
	}

	o.Document = newDoc
	o.Revision = i.CurrentRevision + 1

	return
}
