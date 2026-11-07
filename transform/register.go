package transform

import (
	"github.com/bits-engine/bits-core/inheritance"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/world"
)

var registerLog = logging.New("bits-core:transform/register")

type Types struct {
	GlobalTransform *component.ComponentType[GlobalTransform]
	Transform       *component.ComponentType[Transform]
}

type Systems struct {
	GlobalTransformCompute world.SystemID
}

func Register(w *world.World, inherTypes *inheritance.Types) (*Types, *Systems) {
	defer registerLog.Get().Debug("types and systems registered")
	globalTransformType := component.Register[GlobalTransform](w.CS())
	transformType := component.Register[Transform](w.CS())

	types := Types{
		GlobalTransform: globalTransformType,
		Transform:       transformType,
	}

	gtComputeSys := w.AddSystem(world.ScheduleUpdate, world.NewSysConf(&globalTransformComputeSystem{
		inherTypes: inherTypes,
		types:      &types,
	}))

	return &types, &Systems{
		GlobalTransformCompute: gtComputeSys,
	}
}
