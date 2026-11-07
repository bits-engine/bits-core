package transform

import "github.com/bits-engine/bits-core/gmath"

// Local transformations of entity
type Transform struct {
	Position gmath.Vector3
	Rotation gmath.Quaternion
	Scale    gmath.Vector3
}

func (t *Transform) toGlobal() GlobalTransform {
	return GlobalTransform{
		t: t.copy(),
	}
}

func (t *Transform) copy() Transform {
	return Transform{
		Position: t.Position,
		Rotation: t.Rotation,
		Scale:    t.Scale,
	}
}

func (t *Transform) Combine(child Transform) Transform {
	return Transform{
		Position: t.Position.Add(t.Rotation.RotateVector(child.Position.Mul(t.Scale))),
		Rotation: t.Rotation.Mul(child.Rotation),
		Scale:    t.Scale.Mul(child.Scale),
	}
}

// Global computed transformations of entity, based on inheritance
type GlobalTransform struct {
	t Transform
}

// Position retrieves computed position.
func (gt *GlobalTransform) Position() gmath.Vector3 {
	return gt.t.Position
}

// Pos retrieves computed position.
func (gt *GlobalTransform) Pos() gmath.Vector3 {
	return gt.Position()
}

// Scale retrieves computed scale.
func (gt *GlobalTransform) Scale() gmath.Vector3 {
	return gt.t.Scale
}

// Rotation retrieves computed rotation.
func (gt *GlobalTransform) Rotation() gmath.Quaternion {
	return gt.t.Rotation
}

// Rot retrieves computed rotation.
func (gt *GlobalTransform) Rot() gmath.Quaternion {
	return gt.Rotation()
}
