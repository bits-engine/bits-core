package gmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func approxEqualVec3(t *testing.T, got, want Vector3, msg string) {
	t.Helper()
	assert.InDelta(t, want.X, got.X, testEpsilon, msg+" (X)")
	assert.InDelta(t, want.Y, got.Y, testEpsilon, msg+" (Y)")
	assert.InDelta(t, want.Z, got.Z, testEpsilon, msg+" (Z)")
}

// --- Constructors ---

func TestConstructors3(t *testing.T) {
	tests := []struct {
		name string
		got  Vector3
		want Vector3
	}{
		{"Vec3", Vec3(3, 4, 5), Vector3{X: 3, Y: 4, Z: 5}},
		{"Zero3", Zero3(), Vector3{X: 0, Y: 0, Z: 0}},
		{"One3", One3(), Vector3{X: 1, Y: 1, Z: 1}},
		{"Up3", Up3(), Vector3{X: 0, Y: 1, Z: 0}},
		{"Down3", Down3(), Vector3{X: 0, Y: -1, Z: 0}},
		{"Right3", Right3(), Vector3{X: 1, Y: 0, Z: 0}},
		{"Left3", Left3(), Vector3{X: -1, Y: 0, Z: 0}},
		{"Forward3", Forward3(), Vector3{X: 0, Y: 0, Z: 1}},
		{"Back3", Backward3(), Vector3{X: 0, Y: 0, Z: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// --- Arithmetic ---

func TestAdd3(t *testing.T) {
	got := Vec3(1, 2, 3).Add(Vec3(4, 5, 6))
	approxEqualVec3(t, got, Vec3(5, 7, 9), "Add")
}

func TestSub3(t *testing.T) {
	got := Vec3(5, 7, 9).Sub(Vec3(1, 2, 3))
	approxEqualVec3(t, got, Vec3(4, 5, 6), "Sub")
}

func TestScale3(t *testing.T) {
	approxEqualVec3(t, Vec3(2, 3, 4).Scale(2), Vec3(4, 6, 8), "Scale positive")
	approxEqualVec3(t, Vec3(2, 3, 4).Scale(-1), Vec3(-2, -3, -4), "Scale negative")
	approxEqualVec3(t, Vec3(2, 3, 4).Scale(0), Zero3(), "Scale zero")
}

func TestMul3(t *testing.T) {
	got := Vec3(2, 3, 4).Mul(Vec3(5, 6, 7))
	approxEqualVec3(t, got, Vec3(10, 18, 28), "Mul")
}

func TestNeg3(t *testing.T) {
	got := Vec3(2, -3, 4).Neg()
	approxEqualVec3(t, got, Vec3(-2, 3, -4), "Neg")
}

func TestAbs3(t *testing.T) {
	tests := []struct {
		in   Vector3
		want Vector3
	}{
		{Vec3(-2, -3, -4), Vec3(2, 3, 4)},
		{Vec3(2, -3, 4), Vec3(2, 3, 4)},
		{Vec3(0, 0, 0), Vec3(0, 0, 0)},
	}
	for _, tt := range tests {
		approxEqualVec3(t, tt.in.Abs(), tt.want, "Abs")
	}
}

// --- Length / distance ---

func TestLength3(t *testing.T) {
	// вектор (2, 3, 6) имеет длину 7 (2²+3²+6² = 4+9+36 = 49)
	assert.InDelta(t, 7, Vec3(2, 3, 6).Length(), testEpsilon, "Length")
	assert.InDelta(t, 0, Zero3().Length(), testEpsilon, "Length of zero vector")
}

func TestLengthSquared3(t *testing.T) {
	assert.InDelta(t, 49, Vec3(2, 3, 6).LengthSquared(), testEpsilon, "LengthSquared")
}

func TestDistance3(t *testing.T) {
	assert.InDelta(t, 7, Vec3(0, 0, 0).Distance(Vec3(2, 3, 6)), testEpsilon, "Distance")
	assert.InDelta(t, 7, Vec3(2, 3, 6).Distance(Vec3(0, 0, 0)), testEpsilon, "Distance symmetric")
}

func TestDistanceSquared3(t *testing.T) {
	assert.InDelta(t, 49, Vec3(0, 0, 0).DistanceSquared(Vec3(2, 3, 6)), testEpsilon, "DistanceSquared")
}

// --- Normalize ---

func TestNormalize3(t *testing.T) {
	got := Vec3(2, 3, 6).Normalize()
	approxEqualVec3(t, got, Vec3(2.0/7, 3.0/7, 6.0/7), "Normalize")
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Normalized length should be 1")
}

func TestNormalizeZeroVector3(t *testing.T) {
	got := Zero3().Normalize()
	approxEqualVec3(t, got, Zero3(), "Normalize of zero vector should return zero, not NaN")
	assert.False(t, math.IsNaN(float64(got.X)), "Normalize should not produce NaN (X)")
	assert.False(t, math.IsNaN(float64(got.Y)), "Normalize should not produce NaN (Y)")
	assert.False(t, math.IsNaN(float64(got.Z)), "Normalize should not produce NaN (Z)")
}

func TestNormalizePreservesDirection3(t *testing.T) {
	got := Vec3(10, 0, 0).Normalize()
	approxEqualVec3(t, got, Vec3(1, 0, 0), "Normalize direction")
}

// --- Dot / Cross ---

func TestDot3(t *testing.T) {
	tests := []struct {
		name string
		a, b Vector3
		want float32
	}{
		{"perpendicular", Vec3(1, 0, 0), Vec3(0, 1, 0), 0},
		{"parallel same direction", Vec3(2, 0, 0), Vec3(3, 0, 0), 6},
		{"opposite direction", Vec3(1, 0, 0), Vec3(-1, 0, 0), -1},
		{"general", Vec3(1, 2, 3), Vec3(4, 5, 6), 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, tt.a.Dot(tt.b), testEpsilon)
		})
	}
}

func TestCross3(t *testing.T) {
	// стандартные базисные соотношения правой руки: X × Y = Z, Y × Z = X, Z × X = Y
	approxEqualVec3(t, Right3().Cross(Up3()), Forward3(), "X cross Y should equal Z")
	approxEqualVec3(t, Up3().Cross(Forward3()), Right3(), "Y cross Z should equal X")
	approxEqualVec3(t, Forward3().Cross(Right3()), Up3(), "Z cross X should equal Y")

	// кросс коллинеарных векторов — нулевой вектор
	approxEqualVec3(t, Vec3(2, 0, 0).Cross(Vec3(4, 0, 0)), Zero3(), "Cross of collinear vectors is zero")

	// результат кросса перпендикулярен обоим исходным векторам
	a := Vec3(1, 2, 3)
	b := Vec3(4, 5, 6)
	cross := a.Cross(b)
	assert.InDelta(t, 0, a.Dot(cross), testEpsilon, "Cross result should be perpendicular to a")
	assert.InDelta(t, 0, b.Dot(cross), testEpsilon, "Cross result should be perpendicular to b")
}

func TestAngleBetween3(t *testing.T) {
	assert.InDelta(t, math.Pi/2, Vec3(1, 0, 0).AngleBetween(Vec3(0, 1, 0)), testEpsilon, "AngleBetween 90 degrees")
	assert.InDelta(t, 0, Vec3(1, 0, 0).AngleBetween(Vec3(1, 0, 0)), testEpsilon, "AngleBetween same vector")
	assert.InDelta(t, math.Pi, Vec3(1, 0, 0).AngleBetween(Vec3(-1, 0, 0)), testEpsilon, "AngleBetween opposite vectors")
}

func TestAngleBetweenZeroVector3(t *testing.T) {
	// защита от деления на ноль, когда один из векторов нулевой
	got := Zero3().AngleBetween(Vec3(1, 0, 0))
	assert.False(t, math.IsNaN(float64(got)), "AngleBetween with zero vector should not produce NaN")
}

// --- Rotation ---

func TestRotateX3(t *testing.T) {
	got := Vec3(0, 1, 0).RotateX(math.Pi / 2)
	approxEqualVec3(t, got, Vec3(0, 0, 1), "RotateX 90 degrees")
}

func TestRotateY3(t *testing.T) {
	got := Vec3(0, 0, 1).RotateY(math.Pi / 2)
	approxEqualVec3(t, got, Vec3(1, 0, 0), "RotateY 90 degrees")
}

func TestRotateZ3(t *testing.T) {
	got := Vec3(1, 0, 0).RotateZ(math.Pi / 2)
	approxEqualVec3(t, got, Vec3(0, 1, 0), "RotateZ 90 degrees")
}

func TestRotateAxis3(t *testing.T) {
	// поворот вокруг оси Y на 90° должен совпадать с RotateY
	v := Vec3(0, 0, 1)
	got := v.RotateAxis(Up3(), math.Pi/2)
	approxEqualVec3(t, got, Vec3(1, 0, 0), "RotateAxis around Y should match RotateY")

	// поворот вектора вокруг самого себя (коллинеарной оси) не должен его менять
	v = Vec3(3, 4, 0).Normalize()
	got = v.RotateAxis(v, 1.2345)
	approxEqualVec3(t, got, v, "RotateAxis around itself should be identity")
}

func TestRotatePreservesLength3(t *testing.T) {
	v := Vec3(3, 4, 5)
	tests := []struct {
		name string
		got  Vector3
	}{
		{"RotateX", v.RotateX(1.234)},
		{"RotateY", v.RotateY(1.234)},
		{"RotateZ", v.RotateZ(1.234)},
		{"RotateAxis", v.RotateAxis(Vec3(1, 1, 1).Normalize(), 1.234)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, v.Length(), tt.got.Length(), testEpsilon, "rotation should preserve length")
		})
	}
}

// --- Interpolation ---

func TestLerp3(t *testing.T) {
	a := Vec3(0, 0, 0)
	b := Vec3(10, 20, 30)

	approxEqualVec3(t, a.Lerp(b, 0), a, "Lerp t=0 should return start")
	approxEqualVec3(t, a.Lerp(b, 1), b, "Lerp t=1 should return end")
	approxEqualVec3(t, a.Lerp(b, 0.5), Vec3(5, 10, 15), "Lerp t=0.5 should return midpoint")
}

func TestMoveTowards3(t *testing.T) {
	a := Vec3(0, 0, 0)
	b := Vec3(10, 0, 0)

	approxEqualVec3(t, a.MoveTowards(b, 100), b, "MoveTowards should clamp to target when maxDelta exceeds distance")
	approxEqualVec3(t, a.MoveTowards(b, 3), Vec3(3, 0, 0), "MoveTowards should move by maxDelta")
	approxEqualVec3(t, b.MoveTowards(b, 5), b, "MoveTowards should stay when already at target")
}

// --- Utils ---

func TestEquals3(t *testing.T) {
	a := Vec3(1, 2, 3)
	b := Vec3(1, 2, 3)
	assert.True(t, a.Equals(b), "expected %v to equal %v", a, b)

	c := Vec3(1, 2, 3.1)
	assert.False(t, a.Equals(c), "expected %v to not equal %v", a, c)

	d := Vec3(1+Epsilon/2, 2, 3)
	assert.True(t, a.Equals(d), "expected values within epsilon to be equal")
}

func TestIsZero3(t *testing.T) {
	assert.True(t, Zero3().IsZero(), "expected zero vector to report true")
	assert.False(t, Vec3(0.1, 0, 0).IsZero(), "expected non-zero vector to report false")
}

func TestMin3(t *testing.T) {
	got := Vec3(1, 5, 3).Min(Vec3(3, 2, 6))
	approxEqualVec3(t, got, Vec3(1, 2, 3), "Min")
}

func TestMax3(t *testing.T) {
	got := Vec3(1, 5, 3).Max(Vec3(3, 2, 6))
	approxEqualVec3(t, got, Vec3(3, 5, 6), "Max")
}

func TestClamp3(t *testing.T) {
	minV := Vec3(0, 0, 0)
	maxV := Vec3(10, 10, 10)

	approxEqualVec3(t, Vec3(5, 5, 5).Clamp(minV, maxV), Vec3(5, 5, 5), "Clamp within range stays unchanged")
	approxEqualVec3(t, Vec3(-5, 15, 5).Clamp(minV, maxV), Vec3(0, 10, 5), "Clamp should clamp out-of-range components")
}

func TestReflect3(t *testing.T) {
	normal := Vec3(0, 1, 0)

	got := Vec3(1, -1, 0).Reflect(normal)
	approxEqualVec3(t, got, Vec3(1, 1, 0), "Reflect off horizontal surface")

	got = Vec3(0, -1, 0).Reflect(normal)
	approxEqualVec3(t, got, Vec3(0, 1, 0), "Reflect straight into normal reverses direction")

	got = Vec3(1, 0, 0).Reflect(normal)
	approxEqualVec3(t, got, Vec3(1, 0, 0), "Reflect parallel to surface is unchanged")
}

func TestString3(t *testing.T) {
	got := Vec3(1, 2, 3).String()
	assert.Equal(t, "Vec3(X=1.000000, Y=2.000000, Z=3.000000)", got)
}

// --- Conversion ---

func TestXY(t *testing.T) {
	got := Vec3(1, 2, 3).XY()
	assert.Equal(t, Vec2(1, 2), got)
}

func TestXZ(t *testing.T) {
	got := Vec3(1, 2, 3).XZ()
	assert.Equal(t, Vec2(1, 3), got)
}

// --- Benchmarks ---

func BenchmarkAdd3(b *testing.B) {
	v1, v2 := Vec3(1, 2, 3), Vec3(4, 5, 6)
	for b.Loop() {
		_ = v1.Add(v2)
	}
}

func BenchmarkCross3(b *testing.B) {
	v1, v2 := Vec3(1, 2, 3), Vec3(4, 5, 6)
	for b.Loop() {
		_ = v1.Cross(v2)
	}
}

func BenchmarkNormalize3(b *testing.B) {
	v := Vec3(3, 4, 5)
	for b.Loop() {
		_ = v.Normalize()
	}
}

func BenchmarkRotateAxis3(b *testing.B) {
	v := Vec3(3, 4, 5)
	axis := Vec3(0, 1, 0)
	for b.Loop() {
		_ = v.RotateAxis(axis, 1.0)
	}
}
