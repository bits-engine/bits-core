package events

import (
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

var log = logging.New("bits-core:events")

type eventUpdateSystem[T any] struct {
	eventResourceType *resource.ResourceType[Events[T]]
}

func (s *eventUpdateSystem[T]) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	cfg.WritesRes(s.eventResourceType)
	return cfg
}

func (s *eventUpdateSystem[T]) Run(w *world.World) {
	res, exists := resource.Get(w.RS(), s.eventResourceType)
	if !exists {
		log.Get().Error("event resource does not exists!", "resourceID", s.eventResourceType.ID())
		panic("event resource does not exists!")
	}

	res.swap()
}
