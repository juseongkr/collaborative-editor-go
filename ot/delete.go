package ot

import (
	"fmt"
	"io"
	"io/ioutil"
)

type Delete struct {
	Count int
}

func (d Delete) InputLength() int {
	return d.Count
}

func (d Delete) OutputLength() int {
	return 0
}

func (d Delete) String() string {
	return fmt.Sprintf("Delete(%d)", d.Count)
}

func (d Delete) Slice(start, end int) PrimitiveOp {
	return Delete{Count: end - start}
}

func (d Delete) Apply(r io.Reader, w io.Writer) error {
	_, err := io.CopyN(ioutil.Discard, r, int64(d.Count))

	return err
}

func (d Delete) Join(next PrimitiveOp) PrimitiveOp {
	if nextDelete, ok := next.(Delete); ok {
		return Delete{Count: d.Count + nextDelete.Count}
	}

	return nil
}

func (d Delete) Swap(next PrimitiveOp) bool {
	if _, ok := next.(Insert); ok {
		return true
	}

	return false
}

func (d Delete) Compose(b PrimitiveOp) PrimitiveOp {
	checkComposeLength(d, b)

	switch b.(type) {
	case NoOp:
		return d
	default:
		panic(ErrUnexpectedOp)
	}
}

func (d Delete) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(d, b)

	switch b := b.(type) {
	case Delete:
		return NoOp{}, NoOp{}
	case Retain:
		return Delete{Count: b.Count}, NoOp{}
	default:
		panic(ErrUnexpectedOp)
	}
}
