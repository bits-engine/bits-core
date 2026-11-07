package transform

import (
	"github.com/bits-engine/bits-core/inheritance"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
	"github.com/bits-engine/bits-ecs/world"
)

type globalTransformComputeSystem struct {
	inherTypes          *inheritance.Types
	types               *Types
	rootEntitiesFilter  component.FilterID
	onlyTransformFilter component.FilterID
}

func (s *globalTransformComputeSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := (&world.AccessConfig{}).Reads(s.inherTypes.Children).
		Reads(s.inherTypes.Children).
		Reads(s.types.Transform).
		Writes(s.types.GlobalTransform)

	s.rootEntitiesFilter = fr.Register(
		cfg,
		component.NewFilter().
			Require(s.inherTypes.Children).
			Exclude(s.inherTypes.Parent),
	)
	s.onlyTransformFilter = fr.Register(
		cfg,
		component.NewFilter().
			Exclude(s.inherTypes.Children).
			Exclude(s.inherTypes.Parent).
			Require(s.types.Transform),
	)

	return cfg
}

func (s *globalTransformComputeSystem) Run(w *world.World) {
	// Take all with children, but without parent
	qr := component.Query(w.CS(), s.rootEntitiesFilter)
	toVisit := []entity.Entity{}

	for qr.Next() {
		toVisit = append(toVisit, qr.Entity())

		if transform, exists := component.Field(qr, s.types.Transform); exists {
			component.Set(w.CS(), qr.Entity(), s.types.GlobalTransform, transform.toGlobal())
		}
	}

	// Iterate over with BFS -> calculating transform
	for i := 0; i < len(toVisit); i++ {
		parent := toVisit[i]
		// Get children
		children, exists := component.Get(w.CS(), parent, s.inherTypes.Children)
		if !exists {
			continue
		}
		// Get globaltransform of parent
		parentGlobalTransform, parentGTExists := component.Get(w.CS(), parent, s.types.GlobalTransform)

		// Iterate over children
		for _, child := range children.List() {
			var newTransform Transform

			toVisit = append(toVisit, child)
			childTransform, childTransformExists := component.Get(w.CS(), child, s.types.Transform)

			if parentGTExists && childTransformExists {
				newTransform = parentGlobalTransform.t.Combine(*childTransform)
			} else if parentGTExists {
				newTransform = parentGlobalTransform.t
			} else if childTransformExists {
				newTransform = *childTransform
			}
			component.Set(w.CS(), child, s.types.GlobalTransform, newTransform.toGlobal())
		}
	}

	// Take all without children and parent, but with Transform
	qr = component.Query(w.CS(), s.onlyTransformFilter)

	for qr.Next() {
		transform, _ := component.Field(qr, s.types.Transform)
		component.Set(w.CS(), qr.Entity(), s.types.GlobalTransform, transform.toGlobal())
	}
}
