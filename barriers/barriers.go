// barriers provides default barrier systems.
//
// Barrier systems it is just noop exclusive systems. They are doing nothing, needed to just separate other systems.
package barriers

import (
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/world"
)

var log = logging.New("bits-core:barriers")

// Barriers contains [world.SystemID]s for barriers.
//
// Barriers order is:
// 
// 1. EventsUpdate
// 
// 1. Physics
// 
// 1. Render
type Barriers struct {
	BeforeEventsUpdate world.SystemID
	AfterEventsUpdate world.SystemID
	BeforePhysics world.SystemID
	AfterPhysics world.SystemID
	BeforeRender world.SystemID
	AfterRender world.SystemID
}
