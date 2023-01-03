// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"image/gif"
	"io"
)

func init() {
	RegisterCodec(&GIFCodec{})
}

// GIFCodec is the GIF codec.
type GIFCodec struct{}

// New returns a new instance of GIFCodec.
func (c *GIFCodec) New() Codec { return &GIFCodec{} }

// Name returns the name of the GIF codec: "gif"
func (c *GIFCodec) Name() string { return "gif" }

// Magic returns magic strings that identify GIF data.
func (c *GIFCodec) Magic() []string {
	return []string{"GIF89a", "GIF87a"}
}

// Decode decodes a GIF image according to the options specified.
func (c *GIFCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return gif.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a GIF image
// without decoding the image.
func (c *GIFCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return gif.DecodeConfig(r)
}

// Encode encodes a GIF image according to the options specified.
func (c *GIFCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	var gifOpt gif.Options

	return gif.Encode(w, i, &gifOpt)
}
