package barriers

import "github.com/bits-engine/bits-ecs/world"

// RegisterBarriers registers barriers to scheduler.
//
// Automatically links up all barriers and added to barriers config systems.
func RegisterBarriers(
	w *world.World,
	cfg *barriersConfig,
) *Barriers {
	sched := w.Sched(world.ScheduleUpdate)

	barriers := &Barriers{}

	beforeEventsUpdateCfg := world.NewSysConf(&exclusiveNoopSystem{})
	afterEventsUpdateCfg := world.NewSysConf(&exclusiveNoopSystem{})
	beforePhysicsCfg := world.NewSysConf(&exclusiveNoopSystem{})
	afterPhysicsCfg := world.NewSysConf(&exclusiveNoopSystem{})
	beforeRenderCfg := world.NewSysConf(&exclusiveNoopSystem{})
	afterRenderCfg := world.NewSysConf(&exclusiveNoopSystem{})

	for _, sysID := range cfg.beforeEventsUpdate {
		beforeEventsUpdateCfg.After(sysID)
	}

	for _, sysID := range cfg.inEventsUpdate {
		beforeEventsUpdateCfg.Before(sysID)
		afterEventsUpdateCfg.After(sysID)
	}

	for _, sysID := range cfg.afterEventsUpdate {
		afterEventsUpdateCfg.Before(sysID)
		beforePhysicsCfg.After(sysID)
	}

	for _, sysID := range cfg.inPhysics {
		beforePhysicsCfg.Before(sysID)
		afterPhysicsCfg.After(sysID)
	}

	for _, sysID := range cfg.afterPhysics {
		afterPhysicsCfg.Before(sysID)
		beforeRenderCfg.After(sysID)
	}

	for _, sysID := range cfg.inRender {
		beforeRenderCfg.Before(sysID)
		afterRenderCfg.After(sysID)
	}

	for _, sysID := range cfg.afterRender {
		afterRenderCfg.Before(sysID)
	}

	barriers.BeforeEventsUpdate = sched.Add(beforeEventsUpdateCfg)
	barriers.AfterEventsUpdate = sched.Add(afterEventsUpdateCfg.After(barriers.BeforeEventsUpdate))

	barriers.BeforePhysics = sched.Add(beforePhysicsCfg.After(barriers.AfterEventsUpdate))
	barriers.AfterPhysics = sched.Add(afterPhysicsCfg.After(barriers.BeforePhysics))

	barriers.BeforeRender = sched.Add(beforeRenderCfg.After(barriers.AfterPhysics))
	barriers.AfterRender = sched.Add(afterRenderCfg.After(barriers.BeforeRender))

	return barriers
}
