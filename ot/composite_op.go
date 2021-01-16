package ot

import (
	"io"
)

type CompositeOp []PrimitiveOp

func NewCompositeOp(ops ...PrimitiveOp) CompositeOp {
	return normalize(ops)
}

func (c CompositeOp) InputLength() (length int) {
	for _, p := range c {
		length += p.InputLength()
	}

	return
}

func (c CompositeOp) OutputLength() (length int) {
	for _, p := range c {
		length += p.OutputLength()
	}

	return
}

func (c CompositeOp) Apply(r io.Reader, w io.Writer) error {
	for _, p := range c {
		if err := p.Apply(r, w); err != nil {
			return err
		}
	}

	return nil
}

func (c CompositeOp) Transform(b CompositeOp) (aPrime, bPrime CompositeOp) {
	slicedA, slicedB := slice(c, b, inputLengthFunc, inputLengthFunc)

	for i := range slicedA {
		aOp, bOp := slicedA[i], slicedB[i]

		aOpPrime, bOpPrime := aOp.Transform(bOp)
		aPrime = append(aPrime, aOpPrime)
		bPrime = append(bPrime, bOpPrime)
	}

	return NewCompositeOp(aPrime...), NewCompositeOp(bPrime...)
}

func (c CompositeOp) Compose(b CompositeOp) CompositeOp {
	slicedA, slicedB := slice(c, b, outputLengthFunc, inputLengthFunc)

	var res []PrimitiveOp
	for i := range slicedA {
		aOp, bOp := slicedA[i], slicedB[i]
		c := aOp.Compose(bOp)
		res = append(res, c)
	}

	return NewCompositeOp(res...)
}

func Compose(ops ...CompositeOp) CompositeOp {
	if l := len(ops); l == 0 {
		return CompositeOp{}
	} else if l < 2 {
		return ops[0]
	}

	op := ops[0]
	for i := 1; i < len(ops); i++ {
		op = op.Compose(ops[i])
	}

	return op
}
