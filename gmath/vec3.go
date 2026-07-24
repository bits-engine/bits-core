package gmath

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func Vec3(x, y, z float32) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

// Arithmetic

func (v Vector3) Add(o Vector3) Vector3 {
	return Vec3(v.X+o.X, v.Y+o.Y, v.Z+o.Z)
}

func (v Vector3) Sub(o Vector3) Vector3 {
	return Vec3(v.X-o.X, v.Y-o.Y, v.Z-o.Z)
}

func (v Vector3) Scale(s float32) Vector3 {
	return Vec3(v.X*s, v.Y*s, v.Z*s)
}

func (v Vector3) Mul(o Vector3) Vector3 {
	return Vec3(v.X*o.X, v.Y*o.Y, v.Z*o.Z)
}

func (v Vector3) Neg() Vector3 {
	return v.Scale(-1.0)
}

func (v Vector3) Abs() Vector3 {
	return Vec3(max(v.X, -v.X), max(v.Y, -v.Y), max(v.Z, -v.Z))
}

// Length

func (v Vector3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSquared())))
}

func (v Vector3) DistanceSquared(o Vector3) float32 {
	return o.Sub(v).LengthSquared()
}

func (v Vector3) Distance(o Vector3) float32 {
	return o.Sub(v).Length()
}

func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if IsEqual(l, 0, NormalEpsilon) {
		return Zero3()
	}
	return Vec3(v.X/l, v.Y/l, v.Z/l)
}

// Dot / Cross

func (v Vector3) Dot(o Vector3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector3) Cross(o Vector3) Vector3 {
	return Vec3(
		v.Y*o.Z-v.Z*o.Y,
		v.Z*o.X-v.X*o.Z,
		v.X*o.Y-v.Y*o.X,
	)
}

func (v Vector3) AngleBetween(o Vector3) float32 {
	denom := v.Length() * o.Length()
	if IsEqual(denom, 0, NormalEpsilon) {
		return 0
	}
	cosTheta := v.Dot(o) / denom
	cosTheta = min(max(cosTheta, -1), 1)
	return float32(math.Acos(float64(cosTheta)))
}

// Rotation

func (v Vector3) RotateAxis(axis Vector3, rad float32) Vector3 {
	s, c := math.Sincos(float64(rad))
	sf, cf := float32(s), float32(c)

	term1 := v.Scale(cf)
	term2 := axis.Cross(v).Scale(sf)
	term3 := axis.Scale(axis.Dot(v) * (1 - cf))

	return term1.Add(term2).Add(term3)
}

func (v Vector3) RotateX(rad float32) Vector3 {
	s, c := math.Sincos(float64(rad))
	sf, cf := float32(s), float32(c)
	return Vec3(
		v.X,
		v.Y*cf-v.Z*sf,
		v.Y*sf+v.Z*cf,
	)
}

func (v Vector3) RotateY(rad float32) Vector3 {
	s, c := math.Sincos(float64(rad))
	sf, cf := float32(s), float32(c)
	return Vec3(
		v.X*cf+v.Z*sf,
		v.Y,
		-v.X*sf+v.Z*cf,
	)
}

func (v Vector3) RotateZ(rad float32) Vector3 {
	s, c := math.Sincos(float64(rad))
	sf, cf := float32(s), float32(c)
	return Vec3(
		v.X*cf-v.Y*sf,
		v.X*sf+v.Y*cf,
		v.Z,
	)
}

// Interpolation

func (v Vector3) Lerp(o Vector3, t float32) Vector3 {
	return Vec3(
		(o.X-v.X)*t+v.X,
		(o.Y-v.Y)*t+v.Y,
		(o.Z-v.Z)*t+v.Z,
	)
}

func (v Vector3) MoveTowards(o Vector3, maxDelta float32) Vector3 {
	diff := o.Sub(v)
	dist := diff.Length()
	if dist <= maxDelta || IsEqual(dist, 0, NormalEpsilon) {
		return o
	}
	return v.Add(diff.Scale(maxDelta / dist))
}

// Utils

func (v Vector3) Equals(o Vector3) bool {
	return IsEqual(v.X, o.X, Epsilon) &&
		IsEqual(v.Y, o.Y, Epsilon) &&
		IsEqual(v.Z, o.Z, Epsilon)
}

func (v Vector3) IsZero() bool {
	return IsEqual(v.X, 0, Epsilon) &&
		IsEqual(v.Y, 0, Epsilon) &&
		IsEqual(v.Z, 0, Epsilon)
}

func (v Vector3) Min(o Vector3) Vector3 {
	return Vec3(min(v.X, o.X), min(v.Y, o.Y), min(v.Z, o.Z))
}

func (v Vector3) Max(o Vector3) Vector3 {
	return Vec3(max(v.X, o.X), max(v.Y, o.Y), max(v.Z, o.Z))
}

func (v Vector3) Clamp(minV, maxV Vector3) Vector3 {
	return Vec3(
		min(max(v.X, minV.X), maxV.X),
		min(max(v.Y, minV.Y), maxV.Y),
		min(max(v.Z, minV.Z), maxV.Z),
	)
}

func (v Vector3) Reflect(normal Vector3) Vector3 {
	d := v.Dot(normal)
	return v.Sub(normal.Scale(2 * d))
}

func (v Vector3) String() string {
	return fmt.Sprintf("Vec3(X=%f, Y=%f, Z=%f)", v.X, v.Y, v.Z)
}

// Conversion

func (v Vector3) XY() Vector2 {
	return Vec2(v.X, v.Y)
}

func (v Vector3) XZ() Vector2 {
	return Vec2(v.X, v.Z)
}
