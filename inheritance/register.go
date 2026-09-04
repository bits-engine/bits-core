package inheritance

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/world"
)

var registerLog = logging.New("bits-core:inheritance/register")

// Types holds every registered component/resource/system for inheritance
type Types struct {
	Children *component.ComponentType[Children]
	Parent   *component.ComponentType[Parent]
}

// Register registrates needed for inheritance components
func Register(w *world.World) *Types {
	defer registerLog.Get().Debug("types registered")
	childrenType := component.Register[Children](w.CS())
	parentType := component.Register[Parent](w.CS())

	return &Types{
		Children: childrenType,
		Parent:   parentType,
	}
}

// SetWriteAccess sets access for system to work with inheritance
func SetWriteAccess(cfg *world.AccessConfig, types *Types) *world.AccessConfig {
	return cfg.
		Writes(types.Children).
		Writes(types.Parent)
}
