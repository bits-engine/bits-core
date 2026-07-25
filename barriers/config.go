package barriers

import "github.com/bits-engine/bits-ecs/world"

type barriersConfig struct {
	beforeEventsUpdate []world.SystemID
	inEventsUpdate []world.SystemID
	afterEventsUpdate []world.SystemID
	inPhysics []world.SystemID
	afterPhysics []world.SystemID
	inRender []world.SystemID
	afterRender []world.SystemID
}

// NewConfig creates new configuration for registering barriers.
//
// With barriers configuration you can link already registered systems to barriers.
func NewConfig() *barriersConfig {
	return &barriersConfig{}
}

// BeforeEventsUpdate adds system before any barriers.
func (c *barriersConfig) BeforeEventsUpdate(sys world.SystemID) *barriersConfig {
	c.beforeEventsUpdate = append(c.beforeEventsUpdate, sys)
	return c
}

// InEventsUpdate adds system inside EventsUpdate barrier.
func (c *barriersConfig) InEventsUpdate(sys world.SystemID) *barriersConfig {
	c.inEventsUpdate = append(c.inEventsUpdate, sys)
	return c
}

// AfterEventsUpdate adds system after EventsUpdate barrier and before Physics.
func (c *barriersConfig) AfterEventsUpdate(sys world.SystemID) *barriersConfig {
	c.afterEventsUpdate = append(c.afterEventsUpdate, sys)
	return c
}

// InPhysics adds system inside Physics barrier.
func (c *barriersConfig) InPhysics(sys world.SystemID) *barriersConfig {
	c.inPhysics = append(c.inPhysics, sys)
	return c
}

// AfterPhysics adds system after Physics barrier and before Render barrier.
func (c *barriersConfig) AfterPhysics(sys world.SystemID) *barriersConfig {
	c.afterPhysics = append(c.afterPhysics, sys)
	return c
}

// InRender adds system inside Render barrier.
func (c *barriersConfig) InRender(sys world.SystemID) *barriersConfig {
	c.inRender = append(c.inRender, sys)
	return c
}

// AfterRender adds system after Render barrier.
//
// It is the last barrier in list.
func (c *barriersConfig) AfterRender(sys world.SystemID) *barriersConfig {
	c.afterRender = append(c.afterRender, sys)
	return c
}
