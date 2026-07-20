package time

import (
	"time"

	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

var log = logging.New("bits-core:time")

type timeUpdateSystem struct {
	timeType *resource.ResourceType[Time]
}

func (s *timeUpdateSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	return (&world.AccessConfig{}).
		WritesRes(s.timeType)
}

func (s *timeUpdateSystem) Run(w *world.World) {
	timeRes, exists := resource.Get(w.RS(), s.timeType)
	if !exists {
		log.Get().Error("time resource does not exists!", "resourceID", s.timeType.ID())
		panic("time resource does not exists!")
	}

	timeRes.setLast(time.Now())
}
