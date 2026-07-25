package barriers

import "github.com/bits-engine/bits-ecs/world"

type exclusiveNoopSystem struct{}

func (s *exclusiveNoopSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	return cfg.Exclusive(true)
}
func (s *exclusiveNoopSystem) Run(w *world.World) {}
