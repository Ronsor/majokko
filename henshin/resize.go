// Copyright 2023 Ronsor Labs. All rights reserved.

package henshin

import (
	"math"

	"golang.org/x/image/draw"
)

// ResizeStrategy is an interpolation strategy for resizing images.
type ResizeStrategy = draw.Interpolator

var (
	BiLinearStrategy = draw.BiLinear
	NearestStrategy = draw.NearestNeighbor
)

// areaFit returns a new width and height x2i and y2i given a
// maximum area and the original width and height x and y.
func areaFit(x, y, area int) (x2i int, y2i int) {
	if x == 0 || y == 0 { return 0, 0 }

	x1 := float64(x)
	y1 := float64(y)
	z2 := float64(area)
	z2_rt := math.Sqrt(z2)

	if x == y { return int(z2_rt), int(z2_rt) }

	if x1 > y1 {
		z2_y1 := z2 * y1
		y2 := z2_rt

		// Solve inequality: x1 * y2^2 <= z2 * y1
		for !(x1 * (y2 * y2) <= z2_y1) {
			y2 += -1.0
		}

		// Use x:y ratio to find x2: x1 / y1 = x2 / y2
		x2 := (x1 / y1) * y2
		return int(x2), int(y2)
	} else {
		z2_x1 := z2 * x1
		x2 := z2_rt

		// Solve inequality: y1 * x2^2 <= z2 * x1
		for !(y1 * (x2 * x2) <= z2_x1) {
			x2 += -1.0
		}

		// Use y:x ratio to find y2: y1 / x1 = y2 / x2
		y2 := (y1 / x1) * x2
		return int(x2), int(y2)
	}
}

