package time

import "time"

// Time stores information about tick timings.
//
// Can be used for physics with Delta()
type Time struct {
	delta time.Duration
	last  time.Time
	start time.Time
}

// Elapsed returns elapsed duration from simulation start (or from time resource registration)
func (t *Time) Elapsed() time.Duration {
	return t.last.Sub(t.start)
}

// Delta gives difference between last scheduler tick and current tick in time.
//
// Every system that reads Delta will have the same value (if run After update system).
func (t *Time) Delta() time.Duration {
	return t.delta
}

func (t *Time) setLast(last time.Time) {
	if last.Before(t.last) {
		log.Get().Error("new last value is before saved")
		panic("new last value is before saved")
	}
	t.delta = last.Sub(t.last)
	t.last = last
}
