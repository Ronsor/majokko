// Copyright 2023 Ronsor Labs. All rights reserved.

package henshin

import (
	"testing"
)

func TestAreaFit(t *testing.T) {
	w, h := areaFit(1920, 1080, 1024*1024)
	if w != 1365 && h != 768 {
		t.Errorf("Expected (1365, 768) but got (%d, %d)", w, h)
	}
}
