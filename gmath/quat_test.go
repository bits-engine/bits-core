package gmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func approxEqualQuat(t *testing.T, got, want Quaternion, msg string) {
	t.Helper()
	assert.InDelta(t, want.X, got.X, testEpsilon, msg+" (X)")
	assert.InDelta(t, want.Y, got.Y, testEpsilon, msg+" (Y)")
	assert.InDelta(t, want.Z, got.Z, testEpsilon, msg+" (Z)")
	assert.InDelta(t, want.W, got.W, testEpsilon, msg+" (W)")
}

// --- Constructors ---

func TestQuatConstructor(t *testing.T) {
	got := Quat(1, 2, 3, 4)
	assert.Equal(t, Quaternion{X: 1, Y: 2, Z: 3, W: 4}, got)
}

func TestIdentity(t *testing.T) {
	got := Identity()
	assert.Equal(t, Quaternion{X: 0, Y: 0, Z: 0, W: 1}, got)
}

func TestFromAxisAngle(t *testing.T) {
	// поворот на 90° вокруг Y
	q := FromAxisAngle(Up3(), math.Pi/2)
	approxEqualQuat(t, q, Quat(0, float32(math.Sin(math.Pi/4)), 0, float32(math.Cos(math.Pi/4))), "FromAxisAngle 90 around Y")

	// нулевой угол должен давать identity
	q = FromAxisAngle(Up3(), 0)
	approxEqualQuat(t, q, Identity(), "FromAxisAngle zero angle should be identity")
}

func TestFromEuler(t *testing.T) {
	// нулевые углы должны давать identity
	got := FromEuler(0, 0, 0)
	approxEqualQuat(t, got, Identity(), "FromEuler zero angles should be identity")
}

// --- Arithmetic ---

func TestQuatAdd(t *testing.T) {
	got := Quat(1, 2, 3, 4).Add(Quat(5, 6, 7, 8))
	approxEqualQuat(t, got, Quat(6, 8, 10, 12), "Add")
}

func TestQuatSub(t *testing.T) {
	got := Quat(5, 6, 7, 8).Sub(Quat(1, 2, 3, 4))
	approxEqualQuat(t, got, Quat(4, 4, 4, 4), "Sub")
}

func TestQuatScale(t *testing.T) {
	got := Quat(1, 2, 3, 4).Scale(2)
	approxEqualQuat(t, got, Quat(2, 4, 6, 8), "Scale")
}

func TestQuatNeg(t *testing.T) {
	got := Quat(1, -2, 3, -4).Neg()
	approxEqualQuat(t, got, Quat(-1, 2, -3, 4), "Neg")
}

func TestMulIdentity(t *testing.T) {
	q := FromAxisAngle(Up3(), math.Pi/3)

	// identity — нейтральный элемент композиции с обеих сторон
	approxEqualQuat(t, q.Mul(Identity()), q, "q * identity should equal q")
	approxEqualQuat(t, Identity().Mul(q), q, "identity * q should equal q")
}

func TestMulComposesRotations(t *testing.T) {
	// поворот на 45° вокруг Y, применённый дважды через Mul,
	// должен быть эквивалентен повороту на 90° вокруг Y
	q45 := FromAxisAngle(Up3(), math.Pi/4)
	q90 := FromAxisAngle(Up3(), math.Pi/2)

	combined := q45.Mul(q45)
	approxEqualQuat(t, combined, q90, "45+45 degree rotations should compose into 90 degrees")
}

func TestMulNotCommutative(t *testing.T) {
	// повороты вокруг разных осей в общем случае не коммутативны
	qx := FromAxisAngle(Right3(), math.Pi/2)
	qy := FromAxisAngle(Up3(), math.Pi/2)

	ab := qx.Mul(qy)
	ba := qy.Mul(qx)

	assert.False(t, ab.Equals(ba), "rotations around different axes should not commute")
}

// --- Length / normalize ---

func TestQuatLength(t *testing.T) {
	got := Quat(0, 0, 0, 1).Length()
	assert.InDelta(t, 1, got, testEpsilon, "Length of identity should be 1")

	got = Quat(1, 2, 2, 0).Length()
	assert.InDelta(t, 3, got, testEpsilon, "Length (1,2,2,0) should be 3")
}

func TestQuatNormalize(t *testing.T) {
	q := Quat(1, 2, 2, 0)
	got := q.Normalize()
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Normalized quaternion should have length 1")
}

func TestQuatNormalizeZero(t *testing.T) {
	got := Quat(0, 0, 0, 0).Normalize()
	approxEqualQuat(t, got, Identity(), "Normalize of zero quaternion should return identity, not NaN")
	assert.False(t, math.IsNaN(float64(got.X)), "Normalize should not produce NaN")
}

// --- Inversion ---

func TestConjugate(t *testing.T) {
	q := Quat(1, 2, 3, 4)
	got := q.Conjugate()
	approxEqualQuat(t, got, Quat(-1, -2, -3, 4), "Conjugate negates X,Y,Z and keeps W")
}

func TestInverseOfNormalized(t *testing.T) {
	q := FromAxisAngle(Up3(), math.Pi/3).Normalize()
	inv := q.Inverse()

	// q * q^-1 должен дать identity
	got := q.Mul(inv)
	approxEqualQuat(t, got, Identity(), "q * inverse(q) should be identity")
}

func TestInverseMatchesConjugateForNormalized(t *testing.T) {
	q := FromAxisAngle(Right3(), math.Pi/4).Normalize()
	approxEqualQuat(t, q.Inverse(), q.Conjugate(), "Inverse should match Conjugate for normalized quaternion")
}

func TestInverseOfZero(t *testing.T) {
	got := Quat(0, 0, 0, 0).Inverse()
	approxEqualQuat(t, got, Identity(), "Inverse of zero quaternion should return identity, not NaN")
}

// --- Applying rotation ---

func TestRotateVectorIdentity(t *testing.T) {
	v := Vec3(1, 2, 3)
	got := Identity().RotateVector(v)
	approxEqualVec3(t, got, v, "Identity rotation should not change vector")
}

func TestRotateVectorMatchesAxisRotation(t *testing.T) {
	tests := []struct {
		name string
		axis Vector3
		v    Vector3
	}{
		{"around Y", Up3(), Forward3()},
		{"around X", Right3(), Up3()},
		{"around Z", Forward3(), Right3()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := FromAxisAngle(tt.axis, math.Pi/2)
			got := q.RotateVector(tt.v)
			want := tt.v.RotateAxis(tt.axis, math.Pi/2)
			approxEqualVec3(t, got, want, "RotateVector should match Vector3.RotateAxis")
		})
	}
}

func TestRotateVectorPreservesLength(t *testing.T) {
	q := FromAxisAngle(Vec3(1, 1, 1).Normalize(), 1.234)
	v := Vec3(3, 4, 5)
	got := q.RotateVector(v)
	assert.InDelta(t, v.Length(), got.Length(), testEpsilon, "Rotation should preserve vector length")
}

// --- Interpolation ---

func TestQuatLerpEndpoints(t *testing.T) {
	a := Identity()
	b := FromAxisAngle(Up3(), math.Pi/2)

	got := a.Lerp(b, 0)
	approxEqualQuat(t, got, a, "Lerp t=0 should return start")

	got = a.Lerp(b, 1)
	approxEqualQuat(t, got, b, "Lerp t=1 should return end")
}

func TestQuatLerpIsNormalized(t *testing.T) {
	a := Identity()
	b := FromAxisAngle(Up3(), math.Pi/2)

	got := a.Lerp(b, 0.5)
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Lerp result should be normalized")
}

func TestSlerpEndpoints(t *testing.T) {
	a := Identity()
	b := FromAxisAngle(Up3(), math.Pi/2)

	got := a.Slerp(b, 0)
	approxEqualQuat(t, got, a, "Slerp t=0 should return start")

	got = a.Slerp(b, 1)
	approxEqualQuat(t, got, b, "Slerp t=1 should return end")
}

func TestSlerpMidpoint(t *testing.T) {
	// slerp на полпути между identity и поворотом на 90° вокруг Y
	// должен дать поворот ровно на 45°
	a := Identity()
	b := FromAxisAngle(Up3(), math.Pi/2)
	want := FromAxisAngle(Up3(), math.Pi/4)

	got := a.Slerp(b, 0.5)
	approxEqualQuat(t, got, want, "Slerp midpoint should be exact half-angle rotation")
}

func TestSlerpIsNormalized(t *testing.T) {
	a := Identity()
	b := FromAxisAngle(Right3(), math.Pi/3)

	got := a.Slerp(b, 0.3)
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Slerp result should be normalized")
}

func TestSlerpNearlyParallelFallsBackToLerp(t *testing.T) {
	// два очень близких кватерниона не должны вызывать деление на почти-ноль (sin(theta0))
	a := FromAxisAngle(Up3(), 0.001)
	b := FromAxisAngle(Up3(), 0.0011)

	got := a.Slerp(b, 0.5)
	assert.False(t, math.IsNaN(float64(got.X)), "Slerp of nearly-parallel quaternions should not produce NaN")
	assert.InDelta(t, 1, got.Length(), testEpsilon, "Slerp result should still be normalized")
}

func TestSlerpShortestPath(t *testing.T) {
	// q и -q представляют одно и то же вращение; slerp должен выбрать
	// кратчайший путь независимо от знака одного из операндов
	a := FromAxisAngle(Up3(), 0.1)
	b := FromAxisAngle(Up3(), 0.2).Neg() // тот же поворот, инвертированный знак

	got := a.Slerp(b, 0.5)
	want := FromAxisAngle(Up3(), 0.15)

	// сравниваем через RotateVector, т.к. got может отличаться знаком от want
	// при этом представлять одно и то же вращение
	v := Vec3(1, 2, 3)
	approxEqualVec3(t, got.RotateVector(v), want.RotateVector(v), "Slerp should take shortest path regardless of sign")
}

// --- Conversion ---

func TestToAxisAngleRoundTrip(t *testing.T) {
	axis := Vec3(1, 1, 0).Normalize()
	angle := float32(math.Pi / 3)

	q := FromAxisAngle(axis, angle)
	gotAxis, gotAngle := q.ToAxisAngle()

	approxEqualVec3(t, gotAxis, axis, "ToAxisAngle axis round-trip")
	assert.InDelta(t, angle, gotAngle, testEpsilon, "ToAxisAngle angle round-trip")
}

func TestToAxisAngleZeroRotation(t *testing.T) {
	// при нулевом угле ось произвольна, но не должна быть NaN
	_, angle := Identity().ToAxisAngle()
	assert.InDelta(t, 0, angle, testEpsilon, "ToAxisAngle of identity should have zero angle")
}

func TestToEulerRoundTrip(t *testing.T) {
	tests := []struct {
		name                  string
		pitch, yaw, roll float32
	}{
		{"small angles", 0.3, 0.5, 0.2},
		{"pitch only", 0.7, 0, 0},
		{"yaw only", 0, 0.9, 0},
		{"roll only", 0, 0, 1.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := FromEuler(tt.pitch, tt.yaw, tt.roll)
			got := q.ToEuler()

			assert.InDelta(t, tt.pitch, got.X, testEpsilon, "pitch round-trip")
			assert.InDelta(t, tt.yaw, got.Y, testEpsilon, "yaw round-trip")
			assert.InDelta(t, tt.roll, got.Z, testEpsilon, "roll round-trip")
		})
	}
}

// --- Utils ---

func TestQuatEquals(t *testing.T) {
	a := Quat(1, 2, 3, 4)
	b := Quat(1, 2, 3, 4)
	assert.True(t, a.Equals(b), "expected %v to equal %v", a, b)

	c := Quat(1, 2, 3, 4.1)
	assert.False(t, a.Equals(c), "expected %v to not equal %v", a, c)
}

func TestIsIdentity(t *testing.T) {
	assert.True(t, Identity().IsIdentity(), "expected Identity() to report true")
	assert.False(t, Quat(0, 0, 0, 0.9).IsIdentity(), "expected non-identity quaternion to report false")
}

func TestQuatString(t *testing.T) {
	got := Quat(1, 2, 3, 4).String()
	assert.Equal(t, "Quat(X=1.000000, Y=2.000000, Z=3.000000, W=4.000000)", got)
}

// --- Benchmarks ---

func BenchmarkQuatMul(b *testing.B) {
	q1 := FromAxisAngle(Up3(), 0.5)
	q2 := FromAxisAngle(Right3(), 0.7)
	for i := 0; i < b.N; i++ {
		_ = q1.Mul(q2)
	}
}

func BenchmarkRotateVector(b *testing.B) {
	q := FromAxisAngle(Up3(), 0.5)
	v := Vec3(1, 2, 3)
	for i := 0; i < b.N; i++ {
		_ = q.RotateVector(v)
	}
}

func BenchmarkSlerp(b *testing.B) {
	a := Identity()
	c := FromAxisAngle(Up3(), math.Pi/2)
	for i := 0; i < b.N; i++ {
		_ = a.Slerp(c, 0.5)
	}
}

func BenchmarkQuatNormalize(b *testing.B) {
	q := Quat(1, 2, 3, 4)
	for i := 0; i < b.N; i++ {
		_ = q.Normalize()
	}
}
