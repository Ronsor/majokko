// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"

	"golang.org/x/image/webp"
)

func init() {
	RegisterCodec(&WEBPCodec{})
}

// WEBPCodec is the WEBP codec.
type WEBPCodec struct{}

// New returns a new instance of WEBPCodec.
func (c *WEBPCodec) New() Codec { return &WEBPCodec{} }

// Name returns the name of the WEBP codec: "webp"
func (c *WEBPCodec) Name() string { return "webp" }

// Magic returns magic strings that identify WEBP data.
func (c *WEBPCodec) Magic() []string {
	return []string{"RIFF????WEBPVP8"}
}

// Decode decodes a WEBP image according to the options specified.
func (c *WEBPCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return webp.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a WEBP image
// without decoding the image.
func (c *WEBPCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return webp.DecodeConfig(r)
}

var (
	_ Decoder = &WEBPCodec{}
)
