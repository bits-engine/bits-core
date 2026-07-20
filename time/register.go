package time

import (
	"time"

	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

// RegisterTime registers time resource and update system in world.
func RegisterTime(w *world.World) (timeType *resource.ResourceType[Time], timeUpdateSystemType world.SystemID) {
	typ := resource.Register[Time](w.RS())
	now := time.Now()
	resource.Set(w.RS(), typ, Time{start: now, last: now})
	systemType := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&timeUpdateSystem{
		timeType: typ,
	}))

	log.Get().Debug("time resource registered",
		"resourceID", typ.ID(),
		"systemID", systemType,
	)

	return typ, systemType
}
