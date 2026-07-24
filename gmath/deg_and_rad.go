package gmath

import "math"

const deg2rad float32 = math.Pi / 180
const rad2deg float32 = 180 / math.Pi

func ToRad(degrees float32) float32 {
	return degrees * deg2rad
}

func ToDeg(radians float32) float32 {
	return radians * rad2deg
}
