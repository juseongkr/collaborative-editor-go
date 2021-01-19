package ot_test

import (
	"github.com/juseongkr/collaborative-editor-go/ot"
	"math/rand"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func randomPrimitive() ot.PrimitiveOp {
	switch rand.Intn(3) {
	case 0:
		return ot.Retain{Count: 1}
	case 1:
		return ot.Insert{Text: randString(1)}
	case 2:
		return ot.Delete{Count: 1}
	}

	return nil
}

func randomPrimitiveOps(inLength, outLength int) []ot.PrimitiveOp {
	ops := ot.CompositeOp([]ot.PrimitiveOp{})
	for {
		if ops.InputLength() == inLength && ops.OutputLength() == outLength {
			return []ot.PrimitiveOp(ops)
		}

		if ops.InputLength() == inLength {
			ops = append(ops, ot.Insert{Text: randString(1)})
		} else if ops.OutputLength() == outLength {
			ops = append(ops, ot.Delete{Count: 1})
		} else {
			ops = append(ops, randomPrimitive())
		}
	}
}

func randomCompositeOp(inLength, outLength int) ot.CompositeOp {
	return ot.NewCompositeOp(randomPrimitiveOps(inLength, outLength)...)
}

func testEquality(t *testing.T, op1, op2 ot.Op) {
	if op1.InputLength() != op2.InputLength() {
		t.Error("length mismatch")
		return
	}

	input := randString(op1.InputLength())
	op1Output, err := ot.ApplyString(op1, input)
	if err != nil {
		t.Error(err)
		return
	}

	op2Output, err := ot.ApplyString(op2, input)
	if err != nil {
		t.Error(err)
		return
	}

	if op1Output != op2Output {
		t.Errorf("op1(%s) != op2(%s)", op1Output, op2Output)
	}
}
