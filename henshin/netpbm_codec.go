// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"
	"strings"

	"github.com/spakin/netpbm"
)

func init() {
	RegisterCodec(&NetPBMCodec{netpbm.PNM})
	RegisterCodec(&NetPBMCodec{netpbm.PBM})
	RegisterCodec(&NetPBMCodec{netpbm.PGM})
	RegisterCodec(&NetPBMCodec{netpbm.PPM})
	RegisterCodec(&NetPBMCodec{netpbm.PAM})
}

// NetPBMCodec is the NetPBM codec.
type NetPBMCodec struct {
	Format netpbm.Format
}

// New returns a new instance of NetPBMCodec.
func (c *NetPBMCodec) New() Codec {
	return &NetPBMCodec{Format: c.Format}
}

// Name returns the name of the NetPBM codec: "netpbm"
func (c *NetPBMCodec) Name() string { return strings.ToLower(c.Format.String()) }

// Magic returns magic strings that identify NetPBM data.
func (c *NetPBMCodec) Magic() []string {
	return []string{"P1", "P2", "P3", "P4", "P5", "P6", "P7"}
}

// Decode decodes a NetPBM image according to the options specified.
func (c *NetPBMCodec) Decode(r io.Reader, d *DecodeOptions) (image.Image, error) {
	if d == nil { d = DefaultDecodeOptions() }

	netpbmOpt := &netpbm.DecodeOptions{
		Target: c.Format,
	}

	img, cm, err := netpbm.DecodeWithComments(r, netpbmOpt)
	if err != nil { return img, err }

	if d.Metadata != nil {
		d.Metadata.Comments = cm
	}

	return img, err
}

// DecodeConfig returns the color model and dimensions of a NetPBM image
// without decoding the image.
func (c *NetPBMCodec) DecodeConfig(r io.Reader, d *DecodeOptions) (image.Config, error) {
	return netpbm.DecodeConfig(r)
}

// Encode encodes a NetPBM image according to the options specified.
func (c *NetPBMCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	if o == nil { o = DefaultEncodeOptions() }

	var comments []string
	if o.Metadata != nil {
		comments = o.Metadata.Comments
	}

	netpbmOpt := &netpbm.EncodeOptions{
		Format: c.Format,
		Comments: comments,
	}

	return netpbm.Encode(w, i, netpbmOpt)
}

var (
	_ Decoder = &NetPBMCodec{}
	_ Encoder = &NetPBMCodec{}
)
