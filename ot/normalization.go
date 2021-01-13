package ot

type joinable interface {
	Join(PrimitiveOp) PrimitiveOp
}

type swappable interface {
	Swap(PrimitiveOp) bool
}

func joinOps(ops []PrimitiveOp) []PrimitiveOp {
	if len(ops) == 0 {
		return ops
	}

	var resultIdx int
	for i := 1; i < len(ops); i++ {
		if joinableOp, ok := ops[resultIdx].(joinable); ok {
			if joined := joinableOp.Join(ops[i]); joined != nil {
				ops[resultIdx] = joined
				continue
			}
		}

		resultIdx++
		ops[resultIdx] = ops[i]
	}

	return ops[:resultIdx+1]
}

func swapOps(ops []PrimitiveOp) []PrimitiveOp {
	for i := 0; i < len(ops); i++ {
		if swappableOp, ok := ops[i].(swappable); ok {
			if swappableOp.Swap(ops[i+1]) {
				a, b := ops[i], ops[i+1]
				ops[i] = b
				ops[i+1] = a
				i++
			}
		}
	}

	return ops
}

func removeNoOps(ops []PrimitiveOp) []PrimitiveOp {
	i := 0
	for _, o := range ops {
		if _, ok := o.(NoOp); !ok {
			ops[i] = o
			i++
		}
	}
	ops = ops[:i]

	return ops
}

func normalize(ops CompositeOp) CompositeOp {
	ops = removeNoOps(ops)

	for {
		var joinNeeded, swapNeeded bool
		for i := 0; i < len(ops)-1; i++ {
			a, b := ops[i], ops[i+1]
			if joinable, ok := a.(joinable); ok {
				if joinable.Join(b) != nil {
					joinNeeded = true
					break
				}
			}

			if swappable, ok := a.(swappable); ok {
				if swappable.Swap(b) {
					swapNeeded = true
					break
				}
			}
		}

		if joinNeeded {
			ops = joinOps(ops)
		} else if swapNeeded {
			ops = swapOps(ops)
		} else {
			break
		}
	}

	return ops
}
