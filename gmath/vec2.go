package gmath

import (
	"fmt"
	"math"
)

// Vector with 2 components
//
// Right - +X
// Up - +Y
type Vector2 struct {
	X float32
	Y float32
}

func Vec2(x float32, y float32) Vector2 {
	return Vector2{X: x, Y: y}
}

func (v Vector2) Add(o Vector2) Vector2 {
	return Vec2(v.X+o.X, v.Y+o.Y)
}

func (v Vector2) Scale(s float32) Vector2 {
	return Vec2(v.X*s, v.Y*s)
}

func (v Vector2) Sub(o Vector2) Vector2 {
	return Vec2(v.X-o.X, v.Y-o.Y)
}

func (v Vector2) Mul(o Vector2) Vector2 {
	return Vec2(v.X*o.X, v.Y*o.Y)
}

func (v Vector2) Abs() Vector2 {
	return Vec2(max(v.X, -v.X), max(v.Y, -v.Y))
}

func (v Vector2) Neg() Vector2 {
	return v.Scale(-1.0)
}

func (v Vector2) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSquared())))
}

func (v Vector2) DistanceSquared(o Vector2) float32 {
	return o.Sub(v).LengthSquared()
}

func (v Vector2) Distance(o Vector2) float32 {
	return o.Sub(v).Length()
}

func (v Vector2) Normalize() Vector2 {
	l := v.Length()
	if IsEqual(l, 0, NormalEpsilon) {
		return Zero()
	}

	return Vec2(v.X/l, v.Y/l)
}

func (v Vector2) Dot(o Vector2) float32 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vector2) Cross(o Vector2) float32 {
	return v.X*o.Y - v.Y*o.X
}

func (v Vector2) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

func (v Vector2) Rotate(rad float32) Vector2 {
	s, c := math.Sincos(float64(rad))
	return Vec2(
		v.X*float32(c)-v.Y*float32(s),
		v.X*float32(s)+v.Y*float32(c),
	)
}

func (v Vector2) AngleBetween(o Vector2) float32 {
	return float32(math.Atan2(float64(v.Cross(o)), float64(v.Dot(o))))
}

func (v Vector2) Perpendicular() Vector2 {
	return Vec2(-v.Y, v.X)
}

// Interpolation

func (v Vector2) Lerp(o Vector2, t float32) Vector2 {
	return Vec2((o.X-v.X)*t+v.X, (o.Y-v.Y)*t+v.Y)
}

func (v Vector2) MoveTowards(o Vector2, maxDelta float32) Vector2 {
	diff := o.Sub(v)
	dist := diff.Length()

	if dist <= maxDelta || IsEqual(dist, 0, NormalEpsilon) {
		return o
	}

	return v.Add(diff.Scale(maxDelta / dist))
}

// Utils

func (v Vector2) Equals(o Vector2) bool {
	return IsEqual(v.X, o.X, Epsilon) && IsEqual(v.Y, o.Y, Epsilon)
}

func (v Vector2) IsZero() bool {
	return IsEqual(v.X, 0, Epsilon) && IsEqual(v.Y, 0, Epsilon)
}

func (v Vector2) Min(o Vector2) Vector2 {
	return Vec2(min(v.X, o.X), min(v.Y, o.Y))
}

func (v Vector2) Max(o Vector2) Vector2 {
	return Vec2(max(v.X, o.X), max(v.Y, o.Y))
}

func (v Vector2) Clamp(minV, maxV Vector2) Vector2 {
	return Vec2(
		min(max(v.X, minV.X), maxV.X),
		min(max(v.Y, minV.Y), maxV.Y),
	)
}

func (v Vector2) Reflect(normal Vector2) Vector2 {
	d := v.Dot(normal)
	return v.Sub(normal.Scale(2 * d))
}

func (v Vector2) String() string {
	return fmt.Sprintf("Vec2(X=%f, Y=%f)", v.X, v.Y)
}
