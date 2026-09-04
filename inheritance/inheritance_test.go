package inheritance_test

import (
	"testing"

	"github.com/bits-engine/bits-core/inheritance"
	"github.com/bits-engine/bits-ecs/entity"
	"github.com/bits-engine/bits-ecs/world"
)

func TestInheritance_AddRemove(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 2})

	inheritanceTypes := inheritance.Register(w)

}

