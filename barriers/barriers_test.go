package barriers

import (
	"testing"

	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
)

type senderSys struct {
	name string
	ch   chan<- string
}

func (s *senderSys) Access(fr world.FilterRegistry) *world.AccessConfig {
	return &world.AccessConfig{}
}
func (s *senderSys) Run(w *world.World) {
	s.ch <- s.name
}

type stopSys struct{}

func (s *stopSys) Access(fr world.FilterRegistry) *world.AccessConfig {
	return &world.AccessConfig{}
}
func (s *stopSys) Run(w *world.World) {
	w.Stop()
}

func TestBarriers_Registration(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 4})

	res := make(chan string, 100)

	sysA := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysA"}))
	sysB := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysB"}))
	sysC := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysC"}))
	sysD := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysD"}))
	sysE := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysE"}))
	sysF := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysF"}))
	sysG := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&senderSys{ch: res, name: "sysG"}))
	sysStop := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&stopSys{}))

	barriers := RegisterBarriers(w, NewConfig().
		BeforeEventsUpdate(sysA).
		InEventsUpdate(sysB).
		AfterEventsUpdate(sysC).
		InPhysics(sysD).
		AfterPhysics(sysE).
		InRender(sysF).
		AfterRender(sysG).
		AfterRender(sysStop),
	)

	w.Run()

	assert.Less(t, barriers.BeforeEventsUpdate, barriers.AfterEventsUpdate)
	assert.Less(t, barriers.AfterEventsUpdate, barriers.BeforePhysics)
	assert.Less(t, barriers.BeforePhysics, barriers.AfterPhysics)
	assert.Less(t, barriers.AfterPhysics, barriers.BeforeRender)
	assert.Less(t, barriers.BeforeRender, barriers.AfterRender)

	close(res)

	arr := make([]string, 0)
	for msg := range res {
		arr = append(arr, msg)
	}

	assert.Equal(t, arr[0], "sysA")
	assert.Equal(t, arr[1], "sysB")
	assert.Equal(t, arr[2], "sysC")
	assert.Equal(t, arr[3], "sysD")
	assert.Equal(t, arr[4], "sysE")
	assert.Equal(t, arr[5], "sysF")
	assert.Equal(t, arr[6], "sysG")
}
