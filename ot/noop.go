package ot

import (
	"io"
)

type NoOp struct {
}

func (n NoOp) InputLength() int {
	return 0
}

func (n NoOp) OutputLength() int {
	return 0
}

func (n NoOp) Slice(start, end int) PrimitiveOp {
	return NoOp{}
}

func (n NoOp) Apply(io.Reader, io.Writer) error {
	return nil
}

func (n NoOp) String() string {
	return "NoOp"
}

func (n NoOp) Compose(b PrimitiveOp) PrimitiveOp {
	switch b := b.(type) {
	case Insert:
		return b
	default:
		panic(ErrUnexpectedOp)
	}
}

func (n NoOp) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	switch b := b.(type) {
	case Insert:
		return Retain{Count: b.OutputLength()}, b
	default:
		panic(ErrUnexpectedOp)
	}
}
