package ot

import (
	"fmt"
	"io"
)

type Insert struct {
	Text string
}

func (i Insert) InputLength() int {
	return 0
}

func (i Insert) OutputLength() int {
	return len(i.Text)
}

func (i Insert) String() string {
	return fmt.Sprintf("Insert(%s)", i.Text)
}

func (i Insert) Slice(start, end int) PrimitiveOp {
	return Insert{Text: i.Text[start:end]}
}

func (i Insert) Apply(r ioReader, w io.Writer) error {
	_, err := io.WriteString(w, i.Text)

	return err
}

func (i Insert) Join(next PrimitiveOp) PrimitiveOp {
	if nextInsert, ok := next.(Insert); ok {
		return Insert{Text: i.Text + nextInsert.Text}
	}

	return nil
}

func (i Insert) Compose(b PrimitiveOp) PrimitiveOp {
	checkComposeLength(i, b)

	switch b.(type) {
	case Delete:
		return NoOp{}
	case Retain:
		return i
	default:
		panic(ErrUnexpectedOp)
	}
}

func (i Insert) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(i, b)

	switch b.(type) {
	case NoOp:
		return i, Retain{Count: i.OutputLength()}
	default:
		panic(ErrUnexpectedOp)
	}
}
