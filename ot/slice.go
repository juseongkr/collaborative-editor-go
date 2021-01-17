package ot

import (
	"fmt"
)

type stack struct {
	current PrimitiveOp
	idx     int
	ops     []PrimitiveOp
}

type lengthFunc func(op Op) int

func inputLengthFunc(op Op) int {
	return op.InputLength()
}

func outputLengthFunc(op Op) int {
	return op.OutputLength()
}

func (s *stack) next() {
	s.idx++
	if s.idx < len(s.ops) {
		s.current = s.ops[s.idx]
	} else {
		s.current = NoOp{}
	}
}

func (s *stack) isEOF() bool {
	return s.idx >= len(s.ops)
}

func (s *stack) String() string {
	idx := s.idx + 1
	if idx >= len(s.ops) {
		idx = len(s.ops)
	}

	return fmt.Sprintf("stack(current: %s, remaining: %v)", s.current, s.ops[idx:])
}

func newStack(ops []PrimitiveOp) *stack {
	if len(ops) == 0 {
		return &stack{
			current: NoOp{},
			idx:     0,
			ops:     ops,
		}
	}

	return &stack{
		current: ops[0],
		idx:     0,
		ops:     ops,
	}
}

func slice(a, b CompositeOp, aLengthFunc, bLengthFunc lengthFunc) (aSliced, bSliced []PrimitiveOp) {

	if aLengthFunc(a) != bLengthFunc(b) {
		panic(ErrLengthMismatch)
	}

	aStack := newStack(a)
	bStack := newStack(b)

	for {
		if aStack.isEOF() && bStack.isEOF() {
			return
		}

		aLength := aLengthFunc(aStack.current)
		bLength := bLengthFunc(bStack.current)

		if aLength == 0 && !aStack.isEOF() {
			aSliced = append(aSliced, aStack.current)
			bSliced = append(bSliced, NoOp{})
			aStack.next()
			continue
		} else if bLength == 0 && !bStack.isEOF() {
			aSliced = append(aSliced, NoOp{})
			bSliced = append(bSliced, bStack.current)
			bStack.next()
			continue
		}

		if aStack.isEOF() || bStack.isEOF() {
			panic(ErrLengthMismatch)
		}

		if aLength == bLength {
			aSliced = append(aSliced, aStack.current)
			bSliced = append(bSliced, bStack.current)

			aStack.next()
			bStack.next()
		} else if aLength > bLength {
			a1 := aStack.current.Slice(0, bLength)
			a2 := aStack.current.Slice(bLength, aLength)

			aSliced = append(aSliced, a1)
			bSliced = append(bSliced, bStack.current)

			aStack.current = a2
			bStack.next()
		} else {
			b1 := bStack.current.Slice(0, aLength)
			b2 := bStack.current.Slice(aLength, bLength)

			aSliced = append(aSliced, aStack.current)
			bSliced = append(bSliced, b1)

			aStack.next()
			bStack.current = b2
		}
	}
}
