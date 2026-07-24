package gmath

import (
	"fmt"
	"math"
)

type Quaternion struct {
	X, Y, Z, W float32
}

func Quat(x, y, z, w float32) Quaternion {
	return Quaternion{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func FromAxisAngle(axis Vector3, rad float32) Quaternion {
	half := rad * 0.5
	s, c := math.Sincos(float64(half))
	sf, cf := float32(s), float32(c)
	return Quat(axis.X*sf, axis.Y*sf, axis.Z*sf, cf)
}

func FromEuler(pitch, yaw, roll float32) Quaternion {
	sp, cp := math.Sincos(float64(pitch * 0.5))
	sy, cy := math.Sincos(float64(yaw * 0.5))
	sr, cr := math.Sincos(float64(roll * 0.5))

	spf, cpf := float32(sp), float32(cp)
	syf, cyf := float32(sy), float32(cy)
	srf, crf := float32(sr), float32(cr)

	return Quat(
		spf*cyf*crf-cpf*syf*srf,
		cpf*syf*crf+spf*cyf*srf,
		cpf*cyf*srf-spf*syf*crf,
		cpf*cyf*crf+spf*syf*srf,
	)
}

// Arithmetic

func (q Quaternion) Add(o Quaternion) Quaternion {
	return Quat(q.X+o.X, q.Y+o.Y, q.Z+o.Z, q.W+o.W)
}

func (q Quaternion) Sub(o Quaternion) Quaternion {
	return Quat(q.X-o.X, q.Y-o.Y, q.Z-o.Z, q.W-o.W)
}

func (q Quaternion) Scale(s float32) Quaternion {
	return Quat(q.X*s, q.Y*s, q.Z*s, q.W*s)
}

func (q Quaternion) Neg() Quaternion {
	return q.Scale(-1.0)
}

func (q Quaternion) Mul(o Quaternion) Quaternion {
	return Quat(
		q.W*o.X+q.X*o.W+q.Y*o.Z-q.Z*o.Y,
		q.W*o.Y-q.X*o.Z+q.Y*o.W+q.Z*o.X,
		q.W*o.Z+q.X*o.Y-q.Y*o.X+q.Z*o.W,
		q.W*o.W-q.X*o.X-q.Y*o.Y-q.Z*o.Z,
	)
}

// Length / normalization

func (q Quaternion) Dot(o Quaternion) float32 {
	return q.X*o.X + q.Y*o.Y + q.Z*o.Z + q.W*o.W
}

func (q Quaternion) LengthSquared() float32 {
	return q.Dot(q)
}

func (q Quaternion) Length() float32 {
	return float32(math.Sqrt(float64(q.LengthSquared())))
}

func (q Quaternion) Normalize() Quaternion {
	l := q.Length()
	if IsEqual(l, 0, NormalEpsilon) {
		return Identity()
	}
	inv := 1.0 / l
	return q.Scale(inv)
}

// Inversion

func (q Quaternion) Conjugate() Quaternion {
	return Quat(-q.X, -q.Y, -q.Z, q.W)
}

func (q Quaternion) Inverse() Quaternion {
	lenSq := q.LengthSquared()
	if IsEqual(lenSq, 0, NormalEpsilon) {
		return Identity()
	}
	return q.Conjugate().Scale(1.0 / lenSq)
}

// Applying rotation

func (q Quaternion) RotateVector(v Vector3) Vector3 {
	qv := Vec3(q.X, q.Y, q.Z)
	t := qv.Cross(v).Scale(2)
	return v.Add(t.Scale(q.W)).Add(qv.Cross(t))
}

// Interpolation

func (q Quaternion) Lerp(o Quaternion, t float32) Quaternion {
	return Quat(
		(o.X-q.X)*t+q.X,
		(o.Y-q.Y)*t+q.Y,
		(o.Z-q.Z)*t+q.Z,
		(o.W-q.W)*t+q.W,
	).Normalize()
}

func (q Quaternion) Slerp(o Quaternion, t float32) Quaternion {
	dot := q.Dot(o)

	if dot < 0 {
		o = o.Neg()
		dot = -dot
	}

	// при почти параллельных кватернионах sin(theta0) стремится к нулю,
	// что дало бы деление на почти-ноль — переключаемся на Lerp
	if dot > 0.9995 {
		return q.Lerp(o, t)
	}

	dot = min(max(dot, -1), 1)
	theta0 := float32(math.Acos(float64(dot)))
	theta := theta0 * t

	sinTheta0 := float32(math.Sin(float64(theta0)))
	sinTheta := float32(math.Sin(float64(theta)))

	s0 := float32(math.Cos(float64(theta))) - dot*sinTheta/sinTheta0
	s1 := sinTheta / sinTheta0

	return q.Scale(s0).Add(o.Scale(s1))
}

// Conversion

func (q Quaternion) ToAxisAngle() (Vector3, float32) {
	angle := 2 * float32(math.Acos(float64(min(max(q.W, -1), 1))))
	s := float32(math.Sqrt(float64(1 - q.W*q.W)))
	if IsEqual(s, 0, NormalEpsilon) {
		return Right3(), angle
	}
	return Vec3(q.X/s, q.Y/s, q.Z/s), angle
}

func (q Quaternion) ToEuler() Vector3 {
	// pitch (X)
	sinp := 2 * (q.W*q.X + q.Y*q.Z)
	cosp := 1 - 2*(q.X*q.X+q.Y*q.Y)
	pitch := float32(math.Atan2(float64(sinp), float64(cosp)))

	// yaw (Y)
	siny := 2 * (q.W*q.Y - q.Z*q.X)
	siny = min(max(siny, -1), 1) // защита от погрешности на границах
	yaw := float32(math.Asin(float64(siny)))

	// roll (Z)
	sinr := 2 * (q.W*q.Z + q.X*q.Y)
	cosr := 1 - 2*(q.Y*q.Y+q.Z*q.Z)
	roll := float32(math.Atan2(float64(sinr), float64(cosr)))

	return Vec3(pitch, yaw, roll)
}

// Utils

func (q Quaternion) Equals(o Quaternion) bool {
	return IsEqual(q.X, o.X, Epsilon) &&
		IsEqual(q.Y, o.Y, Epsilon) &&
		IsEqual(q.Z, o.Z, Epsilon) &&
		IsEqual(q.W, o.W, Epsilon)
}

func (q Quaternion) IsIdentity() bool {
	return q.Equals(Identity())
}

func (q Quaternion) String() string {
	return fmt.Sprintf("Quat(X=%f, Y=%f, Z=%f, W=%f)", q.X, q.Y, q.Z, q.W)
}
