// inheritance provides functionality for buliding hierarchy of entities.
package inheritance

import (
	"slices"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
)

// Add adds one entity to children list for another.
//
// Call SetWriteAccess() inside system's Access method.
func Add(
	cs *component.ComponentStorage,
	types *Types,
	parent entity.Entity,
	child entity.Entity,
) {
	parentComp, exists := component.Get(cs, child, types.Parent)
	if exists {
		Remove(cs, types, parentComp.Entity(), child)
	} 

	component.Set(cs, child, types.Parent, Parent{
		entity: parent,
	})

	childrenComp, exists := component.Get(cs, parent, types.Children)
	if exists {
		if !slices.Contains(childrenComp.entities, child) {
			childrenComp.entities = append(childrenComp.entities, child)
		}
	} else {
		component.Set(cs, parent, types.Children, Children{
			entities: []entity.Entity{child},
		})
	}
}

// Remove adds one entity to children list for another.
//
// Call SetWriteAccess() inside system's Access method.
func Remove(
	cs *component.ComponentStorage,
	types *Types,
	parent entity.Entity,
	child entity.Entity,
) {
	component.Remove(cs, child, types.Parent)

	childrenComp, exists := component.Get(cs, parent, types.Children)
	if exists {
		idx := slices.Index(childrenComp.entities, child)
		if idx > -1 {
			childrenComp.entities = append(childrenComp.entities[:idx], childrenComp.entities[idx+1:]...)
		}

		if len(childrenComp.entities) == 0 {
			component.Remove(cs, parent, types.Children)
		}
	}
}
