package events

import (
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

// RegisterEvents registers events in world, setting resource and adding update system.
func RegisterEvents[T any](w *world.World) (eventType *resource.ResourceType[Events[T]], updateSystemType world.SystemID) {
	typ := resource.Register[Events[T]](w.RS())
	resource.Set(w.RS(), typ, Events[T]{})
	systemType := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&eventUpdateSystem[T]{
		eventResourceType: typ,
	}))

	return typ, systemType
}
