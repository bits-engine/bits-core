package transform

import (
	"testing"

	"github.com/bits-engine/bits-core/gmath"
	"github.com/bits-engine/bits-core/inheritance"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform_Combine(t *testing.T) {
	parent := Transform{
		Position: gmath.Up3(),
		Rotation: gmath.Identity(),
		Scale:    gmath.One3(),
	}
	child := Transform{
		Position: gmath.Backward3(),
		Rotation: gmath.Identity(),
		Scale:    gmath.One3(),
	}

	n := parent.Combine(child)
	assert.True(t, n.Position.Equals(gmath.Backward3().Add(gmath.Up3())))
}

func TestTransform_CombintationSystem(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 2})

	inheritanceTypes := inheritance.Register(w)
	globalTransformType := component.Register[GlobalTransform](w.CS())
	transformType := component.Register[Transform](w.CS())

	types := &Types{
		GlobalTransform: globalTransformType,
		Transform:       transformType,
	}

	sys := &globalTransformComputeSystem{
		inherTypes: inheritanceTypes,
		types:      types,
	}

	w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(sys))

	// Lone entity, no parent/children
	loneTransform := Transform{
		Position: gmath.Vec3(1, 2, 3),
		Rotation: gmath.FromAxisAngle(gmath.Up3(), 0.3),
		Scale:    gmath.Vec3(2, 2, 2),
	}
	loneEnt := w.New()
	component.Set(w.CS(), loneEnt, transformType, loneTransform)

	// Root entity with 2 children on a single layer
	root2Transform := Transform{
		Position: gmath.Vec3(4, 0, 0),
		Rotation: gmath.FromAxisAngle(gmath.Up3(), 0.5),
		Scale:    gmath.Vec3(1, 1, 1),
	}
	child21Transform := Transform{
		Position: gmath.Vec3(1, 0, 0),
		Rotation: gmath.FromAxisAngle(gmath.Right3(), 0.1),
		Scale:    gmath.Vec3(3, 3, 3),
	}
	child22Transform := Transform{
		Position: gmath.Vec3(0, 1, 0),
		Rotation: gmath.FromAxisAngle(gmath.Forward3(), 0.2),
		Scale:    gmath.Vec3(0.5, 0.5, 0.5),
	}

	root2Ent := w.New()
	component.Set(w.CS(), root2Ent, transformType, root2Transform)

	child21Ent := w.New()
	component.Set(w.CS(), child21Ent, transformType, child21Transform)
	inheritance.Set(w.CS(), inheritanceTypes, root2Ent, child21Ent)

	child22Ent := w.New()
	component.Set(w.CS(), child22Ent, transformType, child22Transform)
	inheritance.Set(w.CS(), inheritanceTypes, root2Ent, child22Ent)

	// Root entity with a 2-layer chain: root3 -> child31 -> child311
	root3Transform := Transform{
		Position: gmath.Vec3(0, 0, 5),
		Rotation: gmath.FromAxisAngle(gmath.Forward3(), 0.4),
		Scale:    gmath.Vec3(2, 1, 1),
	}
	child31Transform := Transform{
		Position: gmath.Vec3(0, 2, 0),
		Rotation: gmath.FromAxisAngle(gmath.Right3(), 0.6),
		Scale:    gmath.Vec3(1, 2, 1),
	}
	child311Transform := Transform{
		Position: gmath.Vec3(3, 0, 0),
		Rotation: gmath.FromAxisAngle(gmath.Up3(), 0.7),
		Scale:    gmath.Vec3(1, 1, 2),
	}

	root3Ent := w.New()
	component.Set(w.CS(), root3Ent, transformType, root3Transform)

	child31Ent := w.New()
	component.Set(w.CS(), child31Ent, transformType, child31Transform)
	inheritance.Set(w.CS(), inheritanceTypes, root3Ent, child31Ent)

	child311Ent := w.New()
	component.Set(w.CS(), child311Ent, transformType, child311Transform)
	inheritance.Set(w.CS(), inheritanceTypes, child31Ent, child311Ent)

	sys.Run(w)

	// Lone entity: global transform equals its local transform
	assertGlobalTransform(t, w, types, loneEnt, loneTransform)

	// Root with 2 children on one layer
	assertGlobalTransform(t, w, types, root2Ent, root2Transform)
	assertGlobalTransform(t, w, types, child21Ent, root2Transform.Combine(child21Transform))
	assertGlobalTransform(t, w, types, child22Ent, root2Transform.Combine(child22Transform))

	// Root with a 2-layer chain of children
	assertGlobalTransform(t, w, types, root3Ent, root3Transform)
	expectedChild31 := root3Transform.Combine(child31Transform)
	assertGlobalTransform(t, w, types, child31Ent, expectedChild31)
	assertGlobalTransform(t, w, types, child311Ent, expectedChild31.Combine(child311Transform))
}

func assertGlobalTransform(
	t *testing.T,
	w *world.World,
	types *Types,
	ent entity.Entity,
	expected Transform,
) {
	t.Helper()

	gt, exists := component.Get(w.CS(), ent, types.GlobalTransform)
	require.True(t, exists, "GlobalTransform does not exist on entity")

	assert.True(t, gt.Position().Equals(expected.Position), "position mismatch: got %s, want %s", gt.Position(), expected.Position)
	assert.True(t, gt.Rotation().Equals(expected.Rotation), "rotation mismatch: got %s, want %s", gt.Rotation(), expected.Rotation)
	assert.True(t, gt.Scale().Equals(expected.Scale), "scale mismatch: got %s, want %s", gt.Scale(), expected.Scale)
}
