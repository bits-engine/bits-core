package inheritance

import "github.com/bits-engine/bits-ecs/entity"

// Children holds IDs of child entities
type Children struct {
	entities []entity.Entity
}

// List returns list of child entities
func (c *Children) List() []entity.Entity {
	return c.entities
}

// Children holds ID of parent entity
type Parent struct {
	entity entity.Entity
}

// Entity returns parent entity
func (p *Parent) Entity() entity.Entity {
	return p.entity
}
