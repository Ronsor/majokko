// Copyright 2023 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"

	"golang.org/x/image/tiff"
)

func init() {
	RegisterCodec(&TIFFCodec{})
}

// TIFFCodec is the TIFF codec.
type TIFFCodec struct{}

// New returns a new instance of TIFFCodec.
func (c *TIFFCodec) New() Codec { return &TIFFCodec{} }

// Name returns the name of the TIFF codec: "tiff"
func (c *TIFFCodec) Name() string { return "tiff" }

// Aliases returns alternate names for the TIFF codec
func (c *TIFFCodec) Aliases() []string {
	return []string{"tif"}
}

// Magic returns magic strings that identify TIFF data.
func (c *TIFFCodec) Magic() []string {
	return []string{"II\x2A\x00", "MM\x00\x2A"}
}

// Decode decodes a TIFF image according to the options specified.
func (c *TIFFCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return tiff.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a TIFF image
// without decoding the image.
func (c *TIFFCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return tiff.DecodeConfig(r)
}

// Encode encodes a TIFF image according to the options specified.
func (c *TIFFCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	var tiffOpt tiff.Options
	// TODO: allow setting these options

	return tiff.Encode(w, i, &tiffOpt)
}

var (
	_ Decoder = &TIFFCodec{}
	_ Encoder = &TIFFCodec{}
	_ CodecWithAliases = &TIFFCodec{}
)
