package ot

import (
	"fmt"
	"io"
)

type Retain struct {
	Count int
}

func (r Retain) InputLength() int {
	return r.Count
}

func (r Retain) OutputLength() int {
	return r.Count
}

func (r Retain) String() string {
	return fmt.Sprintf("Retain(%d)", r.Count)
}

func (r Retain) Slice(start, end int) PrimitiveOp {
	return Retain{Count: end - start}
}

func (r Retain) Apply(r io.Reader, w io.Writer) error {
	_, err := io.CopyN(w, reader, int64(r.Count))

	return err
}

func (r Retain) Join(next PrimitiveOp) PrimitiveOp {
	if nextRetain, ok := next.(Retain); ok {
		return Retain{Count: r.Count + nextRetain.Count}
	}

	return nil
}

func (r Retain) Compose(b PrimitiveOp) PrimitiveOp {
	checkComposeLength(r, b)

	switch b := b.(type) {
	case Retain:
		return r
	case Delete:
		return b
	default:
		panic(ErrUnexpectedOp)
	}
}

func (r Retain) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(r, b)

	switch b := b.(type) {
	case Retain:
		return b, b
	case Delete:
		return NoOp{}, b
	default:
		panic(ErrUnexpectedOp)
	}
}
