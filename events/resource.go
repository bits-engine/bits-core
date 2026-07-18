// events provides resource for system communication
//
// Register events with [events.RegisterEvents] function
package events

// Events is an resource that holds events and allow systems for easy communication.
//
// Works on two buffers - one for read, second for write. On one scheduler tick it swaps.
type Events[T any] struct {
	read []T
	write []T
}

// Send adds new event to event queue. Use this only with WriteRes access configuration.
func (e *Events[T]) Send(event T) {
	e.write = append(e.write, event)
}

// Read returns events for reading. Use this only with ReadRes access configuration.
func (e *Events[T]) Read() []T {
	return e.read
}

func (e *Events[T]) swap() {
	tmp := e.read[:0]
	e.read = e.write
	e.write = tmp
}
