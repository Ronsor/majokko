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
	z1 := x1 * y1

	z2 := float64(area)
	z2z1_rt := math.Sqrt(z2/z1)

	x2 := x1 * z2z1_rt
	y2 := y1 * z2z1_rt

	return int(x2), int(y2)
}

