package ot_test

import (
	"fmt"
	"github.com/juseongkr/collaborative-editor-go/ot"
	"math/rand"
	"testing"
)

func TestCompositeOpApplyString(t *testing.T) {

	op := ot.NewCompositeOp(
		ot.Delete{Count: 1},
		ot.Insert{Text: "H"},
		ot.Retain{Count: 4},
		ot.Insert{Text: ","},
		ot.Retain{Count: 1},
		ot.Delete{Count: 1},
		ot.Insert{Text: "W"},
		ot.Retain{Count: 4},
		ot.Insert{Text: "!"},
	)

	applied, err := ot.ApplyString(op, "hello world")
	if err != nil {
		t.Error(err)
		return
	}

	if applied != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", applied)
	}
}

func TestNewCompositeOp(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("new-op-%d", i), func(t *testing.T) {
			t.Parallel()
			inLength := rand.Intn(10) + 5
			outLength := inLength + rand.Intn(10) - 5

			primitives := randomPrimitiveOps(inLength, outLength)
			primitivesCopy := make([]ot.PrimitiveOp, len(primitives))
			copy(primitivesCopy, primitives)

			op := ot.NewCompositeOp(primitivesCopy...)

			for j := 0; j < len(op)-1; j++ {
				a, b := op[j], op[j+1]
				if joinable, ok := a.(ot.Joinable); ok {
					if joinable.Join(b) != nil {
						t.Errorf("joinable CompositeOp found: %v, ops: %v", joinable, op)
						return
					}
				}

				if swappable, ok := a.(ot.Swappable); ok {
					if swappable.Swap(b) {
						t.Errorf("swappable CompositeOp found: %v, ops: %v", swappable, op)
						return
					}
				}
			}

			t.Run("original-simplified-equality", func(t *testing.T) {
				testEquality(t, ot.CompositeOp(primitives), op)
			})
		})
	}
}
