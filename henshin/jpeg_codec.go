// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"image/jpeg"
	"io"
)

func init() {
	RegisterCodec(&JPEGCodec{})
}

// JPEGCodec is the JPEG codec.
type JPEGCodec struct{}

// New returns a new instance of JPEGCodec.
func (c *JPEGCodec) New() Codec { return &JPEGCodec{} }

// Name returns the name of the JPEG codec: "jpeg"
func (c *JPEGCodec) Name() string { return "jpeg" }

// Aliases returns alternate names for the JPEG codec
func (c *JPEGCodec) Aliases() []string {
	return []string{"jpg", "jfif", "jpi"}
}

// Magic returns magic strings that identify JPEG data.
func (c *JPEGCodec) Magic() []string {
	return []string{"\xff\xd8"}
}

// Decode decodes a JPEG image according to the options specified.
func (c *JPEGCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return jpeg.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a JPEG image
// without decoding the image.
func (c *JPEGCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return jpeg.DecodeConfig(r)
}

// Encode encodes a JPEG image according to the options specified.
func (c *JPEGCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	var jpegOpt jpeg.Options
	if o.CompressionLevel < 100 {
		jpegOpt.Quality = 100 - o.CompressionLevel
	} else {
		jpegOpt.Quality = jpeg.DefaultQuality
	}

	return jpeg.Encode(w, i, &jpegOpt)
}

var (
	_ Decoder = &JPEGCodec{}
	_ Encoder = &JPEGCodec{}
	_ CodecWithAliases = &JPEGCodec{}
)
