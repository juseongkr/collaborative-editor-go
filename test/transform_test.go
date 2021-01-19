package ot_test

import (
	"fmt"
	"github.com/juseongkr/collaborative-editor-go/ot"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func testTransform(t *testing.T, input string, op1, op2 ot.CompositeOp) {
	op1Prime, op2Prime := op1.Transform(op2)

	afterOp1, err := ot.ApplyString(op1, input)
	if err != nil {
		t.Error(err)
		return
	}

	afterOp2Prime, err := ot.ApplyString(op2Prime, afterOp1)
	if err != nil {
		t.Error(err)
		return
	}

	afterOp2, err := ot.ApplyString(op2, input)
	if err != nil {
		t.Error(err)
		return
	}

	afterOp1Prime, err := ot.ApplyString(op1Prime, afterOp2)
	if err != nil {
		t.Error(err)
		return
	}

	if afterOp1Prime != afterOp2Prime {
		t.Errorf("afterOp2Prime(%s) != afterOp1Prime(%s)", afterOp2Prime, afterOp1Prime)
	}
}

func TestTransform(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("rand-%d", i), func(t *testing.T) {
			t.Parallel()
			l := rand.Intn(20) + 6
			inputStr := randString(l)

			op1 := randomCompositeOp(l, l+rand.Intn(10)-5)
			op2 := randomCompositeOp(l, l+rand.Intn(10)-5)

			testTransform(t, inputStr, op1, op2)
		})
	}
}
