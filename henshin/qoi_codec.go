// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"

	"github.com/xfmoulet/qoi"
)

func init() {
	RegisterCodec(&QOICodec{})
}

// QOICodec is the QOI codec.
type QOICodec struct{}

// New returns a new instance of QOICodec.
func (c *QOICodec) New() Codec { return &QOICodec{} }

// Name returns the name of the QOI codec: "qoi"
func (c *QOICodec) Name() string { return "qoi" }

// Magic returns magic strings that identify QOI data.
func (c *QOICodec) Magic() []string {
	return []string{"qoif"}
}

// Decode decodes a QOI image according to the options specified.
func (c *QOICodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	return qoi.Decode(r)
}

// DecodeConfig returns the color model and dimensions of a QOI image
// without decoding the image.
func (c *QOICodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return qoi.DecodeConfig(r)
}

// Encode encodes a QOI image according to the options specified.
func (c *QOICodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	return qoi.Encode(w, i)
}

var (
	_ Decoder = &QOICodec{}
	_ Encoder = &QOICodec{}
)
