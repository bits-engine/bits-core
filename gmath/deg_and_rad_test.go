package gmath_test

import (
	"math"
	"testing"

	"github.com/bits-engine/bits-core/gmath"
	"github.com/stretchr/testify/assert"
)

func TestGMath_DegToRad(t *testing.T) {
	var degrees float32 = 54.0
	radians := gmath.ToRad(degrees)
	
	assert.InEpsilon(t, radians, degrees / 180.0 * math.Pi, float64(gmath.Epsilon))

	deg := gmath.ToDeg(radians)

	assert.InEpsilon(t, deg, degrees, float64(gmath.Epsilon))
}
