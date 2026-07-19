package time

import (
	"testing"
	"time"

	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
)

func areTimesClose(t1, t2 time.Time, margin time.Duration) bool {
	return t1.Sub(t2).Abs() <= margin
}

func TestTime_System(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 1})
	timeType, _ := RegisterTime(w)
	tm, _ := resource.Get(w.RS(), timeType)

	assert.True(t, areTimesClose(tm.start, time.Now(), 100 * time.Millisecond), "start is not near time.Now()")
	lastTime := tm.last

	sys := &timeUpdateSystem{timeType: timeType}
	sys.Run(w)

	assert.True(t, tm.Delta() > 0, "delta can not be less than zero")
	assert.True(t, tm.Delta() < time.Millisecond * 50, "delta is greater than 50ms")
	assert.Equal(t, tm.Delta(), tm.last.Sub(lastTime), "delta calculated incorrectly")
}
