package gmath

const Epsilon float32 = 1e-6
const NormalEpsilon float32 = 1e-4
const PhysicsEpsilon float32 = 1e-3

func IsEqual(a, b float32, eps float32) bool {
	diff := a - b
	diff = max(diff, -diff)
	return diff < eps
}
