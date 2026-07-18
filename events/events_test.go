package events

import (
	"testing"

	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
)

func TestEvents_Swap(t *testing.T) {
	events := Events[int]{}

	assert.Len(t, events.Read(), 0)

	for range 20 {
		events.Send(0)
	}

	assert.Len(t, events.read, 0)
	assert.Len(t, events.write, 20)

	events.swap()

	assert.Len(t, events.Read(), 20)
	assert.Len(t, events.write, 0)

	events.Send(1)
	events.Send(1)

	assert.Len(t, events.Read(), 20)
	assert.Len(t, events.write, 2)

	for _, val := range events.Read() {
		assert.Equal(t, val, 0)
	}

	for _, val := range events.write {
		assert.Equal(t, val, 1)
	}

	events.swap()

	assert.Len(t, events.Read(), 2)
	assert.Len(t, events.write, 0)

	for _, val := range events.Read() {
		assert.Equal(t, val, 1)
	}
}

func TestEvents_SystemSwap(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 1})
	eventType, _ := RegisterEvents[int](w)
	events, _ := resource.Get(w.RS(), eventType)

	events.Send(1)

	assert.Len(t, events.read, 0)
	assert.Len(t, events.write, 1)

	sys := &eventUpdateSystem[int]{eventResourceType: eventType}

	sys.Run(w)

	assert.Len(t, events.read, 1)
	assert.Equal(t, events.read[0], 1)
	assert.Len(t, events.write, 0)
}
