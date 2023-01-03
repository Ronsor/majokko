// Copyright 2023 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"

	"golang.org/x/image/bmp"
)

func init() {
	RegisterCodec(&BMPCodec{})
}

// BMPCodec is the BMP codec.
type BMPCodec struct{}

// New returns a new instance of BMPCodec.
func (c *BMPCodec) New() Codec { return &BMPCodec{} }

// Name returns the name of the BMP codec: "bmp"
func (c *BMPCodec) Name() string { return "bmp" }

// Aliases returns alternate names for the BMP codec
func (c *BMPCodec) Aliases() []string {
	return []string{"dib"}
}

// Magic returns magic strings that identify BMP data.
func (c *BMPCodec) Magic() []string {
	return []string{"BM????\x00\x00\x00\x00"}
}

// Decode decodes a BMP image according to the options specified.
func (c *BMPCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return bmp.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a BMP image
// without decoding the image.
func (c *BMPCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return bmp.DecodeConfig(r)
}

// Encode encodes a BMP image according to the options specified.
func (c *BMPCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	return bmp.Encode(w, i)
}

var (
	_ Decoder = &BMPCodec{}
	_ Encoder = &BMPCodec{}
	_ CodecWithAliases = &BMPCodec{}
)
