package gmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testEpsilon = 1e-5

func approxEqualVec(t *testing.T, got, want Vector2, msg string) {
	t.Helper()
	assert.InDelta(t, want.X, got.X, testEpsilon, msg+" (X)")
	assert.InDelta(t, want.Y, got.Y, testEpsilon, msg+" (Y)")
}

// --- Constructors ---

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  Vector2
		want Vector2
	}{
		{"Vec2", Vec2(3, 4), Vector2{X: 3, Y: 4}},
		{"Zero", Zero(), Vector2{X: 0, Y: 0}},
		{"One", One(), Vector2{X: 1, Y: 1}},
		{"NegOne", NegOne(), Vector2{X: -1, Y: -1}},
		{"Up", Up(), Vector2{X: 0, Y: 1}},
		{"Down", Down(), Vector2{X: 0, Y: -1}},
		{"Right", Right(), Vector2{X: 1, Y: 0}},
		{"Left", Left(), Vector2{X: -1, Y: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// --- Arithmetic ---

func TestAdd(t *testing.T) {
	got := Vec2(1, 2).Add(Vec2(3, 4))
	approxEqualVec(t, got, Vec2(4, 6), "Add")
}

func TestSub(t *testing.T) {
	got := Vec2(5, 7).Sub(Vec2(2, 3))
	approxEqualVec(t, got, Vec2(3, 4), "Sub")
}

func TestScale(t *testing.T) {
	approxEqualVec(t, Vec2(2, 3).Scale(2), Vec2(4, 6), "Scale positive")
	approxEqualVec(t, Vec2(2, 3).Scale(-1), Vec2(-2, -3), "Scale negative")
	approxEqualVec(t, Vec2(2, 3).Scale(0), Zero(), "Scale zero")
}

func TestMul(t *testing.T) {
	got := Vec2(2, 3).Mul(Vec2(4, 5))
	approxEqualVec(t, got, Vec2(8, 15), "Mul")
}

func TestNeg(t *testing.T) {
	got := Vec2(2, -3).Neg()
	approxEqualVec(t, got, Vec2(-2, 3), "Neg")
}

func TestAbs(t *testing.T) {
	tests := []struct {
		in   Vector2
		want Vector2
	}{
		{Vec2(-2, -3), Vec2(2, 3)},
		{Vec2(2, -3), Vec2(2, 3)},
		{Vec2(-2, 3), Vec2(2, 3)},
		{Vec2(0, 0), Vec2(0, 0)},
	}
	for _, tt := range tests {
		approxEqualVec(t, tt.in.Abs(), tt.want, "Abs")
	}
}

// --- Length / distance ---

func TestLength(t *testing.T) {
	assert.InDelta(t, 5, Vec2(3, 4).Length(), testEpsilon, "Length (3-4-5 triangle)")
	assert.InDelta(t, 0, Zero().Length(), testEpsilon, "Length of zero vector")
}

func TestLengthSquared(t *testing.T) {
	assert.InDelta(t, 25, Vec2(3, 4).LengthSquared(), testEpsilon, "LengthSquared")
}

func TestDistance(t *testing.T) {
	assert.InDelta(t, 5, Vec2(0, 0).Distance(Vec2(3, 4)), testEpsilon, "Distance")
	assert.InDelta(t, 5, Vec2(3, 4).Distance(Vec2(0, 0)), testEpsilon, "Distance symmetric")
}

func TestDistanceSquared(t *testing.T) {
	assert.InDelta(t, 25, Vec2(0, 0).DistanceSquared(Vec2(3, 4)), testEpsilon, "DistanceSquared")
}

// --- Normalize ---

func TestNormalize(t *testing.T) {
	got := Vec2(3, 4).Normalize()
	approxEqualVec(t, got, Vec2(0.6, 0.8), "Normalize")
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Normalized length should be 1")
}

func TestNormalizeZeroVector(t *testing.T) {
	got := Zero().Normalize()
	approxEqualVec(t, got, Zero(), "Normalize of zero vector should return zero, not NaN")
	assert.False(t, math.IsNaN(float64(got.X)), "Normalize should not produce NaN (X)")
	assert.False(t, math.IsNaN(float64(got.Y)), "Normalize should not produce NaN (Y)")
}

func TestNormalizePreservesDirection(t *testing.T) {
	got := Vec2(10, 0).Normalize()
	approxEqualVec(t, got, Vec2(1, 0), "Normalize direction")
}

// --- Dot / Cross ---

func TestDot(t *testing.T) {
	tests := []struct {
		name string
		a, b Vector2
		want float32
	}{
		{"perpendicular", Vec2(1, 0), Vec2(0, 1), 0},
		{"parallel same direction", Vec2(2, 0), Vec2(3, 0), 6},
		{"opposite direction", Vec2(1, 0), Vec2(-1, 0), -1},
		{"general", Vec2(2, 3), Vec2(4, 5), 23},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, tt.a.Dot(tt.b), testEpsilon)
		})
	}
}

func TestCross(t *testing.T) {
	tests := []struct {
		name string
		a, b Vector2
		want float32
	}{
		{"b to the left of a (CCW)", Vec2(1, 0), Vec2(0, 1), 1},
		{"b to the right of a (CW)", Vec2(1, 0), Vec2(0, -1), -1},
		{"collinear", Vec2(2, 0), Vec2(4, 0), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, tt.a.Cross(tt.b), testEpsilon)
		})
	}
}

// --- Angles / rotation ---

func TestAngle(t *testing.T) {
	tests := []struct {
		name string
		v    Vector2
		want float32
	}{
		{"along +X axis", Vec2(1, 0), 0},
		{"along +Y axis", Vec2(0, 1), math.Pi / 2},
		{"along -X axis", Vec2(-1, 0), math.Pi},
		{"along -Y axis", Vec2(0, -1), -math.Pi / 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, tt.v.Angle(), testEpsilon)
		})
	}
}

func TestRotate(t *testing.T) {
	approxEqualVec(t, Vec2(1, 0).Rotate(math.Pi/2), Vec2(0, 1), "Rotate 90 degrees")
	approxEqualVec(t, Vec2(1, 0).Rotate(math.Pi), Vec2(-1, 0), "Rotate 180 degrees")

	v := Vec2(3, 4)
	approxEqualVec(t, v.Rotate(0), v, "Rotate 0 degrees is identity")
	approxEqualVec(t, v.Rotate(2*math.Pi), v, "Rotate 360 degrees is identity")
}

func TestRotatePreservesLength(t *testing.T) {
	v := Vec2(3, 4)
	got := v.Rotate(1.234)
	assert.InDelta(t, v.Length(), got.Length(), testEpsilon, "Rotate should preserve length")
}

func TestAngleBetween(t *testing.T) {
	assert.InDelta(t, math.Pi/2, Vec2(1, 0).AngleBetween(Vec2(0, 1)), testEpsilon, "AngleBetween 90 degrees")
	assert.InDelta(t, 0, Vec2(1, 0).AngleBetween(Vec2(1, 0)), testEpsilon, "AngleBetween same vector")
	assert.InDelta(t, -math.Pi/2, Vec2(1, 0).AngleBetween(Vec2(0, -1)), testEpsilon, "AngleBetween signed (clockwise)")
}

func TestPerpendicular(t *testing.T) {
	approxEqualVec(t, Vec2(1, 0).Perpendicular(), Vec2(0, 1), "Perpendicular")

	v := Vec2(3, 4)
	perp := v.Perpendicular()
	assert.InDelta(t, 0, v.Dot(perp), testEpsilon, "Perpendicular dot product should be zero")
}

// --- Interpolation ---

func TestLerp(t *testing.T) {
	a := Vec2(0, 0)
	b := Vec2(10, 20)

	approxEqualVec(t, a.Lerp(b, 0), a, "Lerp t=0 should return start")
	approxEqualVec(t, a.Lerp(b, 1), b, "Lerp t=1 should return end")
	approxEqualVec(t, a.Lerp(b, 0.5), Vec2(5, 10), "Lerp t=0.5 should return midpoint")
}

func TestMoveTowards(t *testing.T) {
	a := Vec2(0, 0)
	b := Vec2(10, 0)

	approxEqualVec(t, a.MoveTowards(b, 100), b, "MoveTowards should clamp to target when maxDelta exceeds distance")
	approxEqualVec(t, a.MoveTowards(b, 3), Vec2(3, 0), "MoveTowards should move by maxDelta")
	approxEqualVec(t, b.MoveTowards(b, 5), b, "MoveTowards should stay when already at target")
}

// --- Utils ---

func TestEquals(t *testing.T) {
	a := Vec2(1, 2)
	b := Vec2(1, 2)
	assert.True(t, a.Equals(b), "expected %v to equal %v", a, b)

	c := Vec2(1, 2.1)
	assert.False(t, a.Equals(c), "expected %v to not equal %v", a, c)

	d := Vec2(1+Epsilon/2, 2)
	assert.True(t, a.Equals(d), "expected values within epsilon to be equal")
}

func TestIsZero(t *testing.T) {
	assert.True(t, Zero().IsZero(), "expected zero vector to report true")
	assert.False(t, Vec2(0.1, 0).IsZero(), "expected non-zero vector to report false")
}

func TestMin(t *testing.T) {
	got := Vec2(1, 5).Min(Vec2(3, 2))
	approxEqualVec(t, got, Vec2(1, 2), "Min")
}

func TestMax(t *testing.T) {
	got := Vec2(1, 5).Max(Vec2(3, 2))
	approxEqualVec(t, got, Vec2(3, 5), "Max")
}

func TestClamp(t *testing.T) {
	minV := Vec2(0, 0)
	maxV := Vec2(10, 10)

	approxEqualVec(t, Vec2(5, 5).Clamp(minV, maxV), Vec2(5, 5), "Clamp within range stays unchanged")
	approxEqualVec(t, Vec2(-5, 15).Clamp(minV, maxV), Vec2(0, 10), "Clamp should clamp out-of-range components")
}

func TestReflect(t *testing.T) {
	normal := Vec2(0, 1)

	got := Vec2(1, -1).Reflect(normal)
	approxEqualVec(t, got, Vec2(1, 1), "Reflect off horizontal surface")

	got = Vec2(0, -1).Reflect(normal)
	approxEqualVec(t, got, Vec2(0, 1), "Reflect straight into normal reverses direction")

	got = Vec2(1, 0).Reflect(normal)
	approxEqualVec(t, got, Vec2(1, 0), "Reflect parallel to surface is unchanged")
}

func TestString(t *testing.T) {
	got := Vec2(1, 2).String()
	assert.Equal(t, "Vec2(X=1.000000, Y=2.000000)", got)
}

// --- Benchmarks ---

func BenchmarkAdd(b *testing.B) {
	v1, v2 := Vec2(1, 2), Vec2(3, 4)
	for b.Loop() {
		_ = v1.Add(v2)
	}
}

func BenchmarkNormalize(b *testing.B) {
	v := Vec2(3, 4)
	for b.Loop() {
		_ = v.Normalize()
	}
}

func BenchmarkRotate(b *testing.B) {
	v := Vec2(3, 4)
	for b.Loop() {
		_ = v.Rotate(1.0)
	}
}

func BenchmarkDot(b *testing.B) {
	v1, v2 := Vec2(1, 2), Vec2(3, 4)
	for b.Loop() {
		_ = v1.Dot(v2)
	}
}
