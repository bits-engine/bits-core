package vec2

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X float32
	Y float32
}

func New(x float32, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

func (v Vec2) Add(o Vec2) Vec2 {
	return New(v.X+o.X, v.Y+o.Y)
}

func (v Vec2) Scale(s float32) Vec2 {
	return New(v.X*s, v.Y*s)
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return New(v.X-o.X, v.Y-o.Y)
}

func (v Vec2) Mul(o Vec2) Vec2 {
	return New(v.X*o.X, v.Y*o.Y)
}

func (v Vec2) Abs(o Vec2) Vec2 {
	return New(max(o.X, -o.X), max(o.Y, -o.Y))
}

func (v Vec2) Neg() Vec2 {
	return v.Scale(-1.0)
}

func (v Vec2) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSquared())))
}

func (v Vec2) DistanceSquared(o Vec2) float32 {
	return o.Sub(v).LengthSquared()
}

func (v Vec2) Distance(o Vec2) float32 {
	return o.Sub(v).Length()
}

func (v Vec2) Normalize() Vec2 {
	if v.LengthSquared() == 0 {
		return Zero()
	}

	l := v.Length()
	return New(v.X/l, v.Y/l)
}

func (v Vec2) String() string {
	return fmt.Sprintf("Vec2(X=%f, Y=%f)", v.X, v.Y)
}
