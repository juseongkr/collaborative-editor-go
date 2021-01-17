package ot

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

var ErrLengthMismatch = errors.New("Length mismatch")
var ErrUnexpectedOp = errors.New("Unexpected operation")

type Op interface {
	InputLength() int
	OutputLength() int
	Apply(io.Reader, io.Writer) error
}

type PrimitiveOp interface {
	Op
	Slice(start, end int) PrimitiveOp
	Compose(b PrimitiveOp) PrimitiveOp
	Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp)
}

func checkComposeLength(a, b Op) {
	if a.OutputLength() != b.InputLength() {
		panic(ErrLengthMismatch)
	}
}

func checkTransformLength(a, b Op) {
	if a.InputLength() != b.InputLength() {
		panic(ErrLengthMismatch)
	}
}

func ApplyString(op Op, text string) (string, error) {
	reader := bytes.NewReader([]byte(text))
	out := new(strings.Builder)
	if err := op.Apply(reader, out); err != nil {
		return "", err
	}

	return out.String(), nil
}
