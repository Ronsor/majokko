// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
)

// diffHash implements the difference-based perceptual hash
// algorithm. It REQUIRES an image.Gray with dimensions of
// 9x8.
func diffHash(i *image.Gray) (ret uint64) {
	size := i.Bounds().Size()
	if size.X != 9 || size.Y != 8 {
		panic("diffHash requires an image with a size of 9x8")
	}

	off := uint64(0)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if i.Pix[y * 9 + x] > i.Pix[y * 9 + x + 1] {
				ret |= 1 << off
			}
			off++
		}
	}
	return
}
