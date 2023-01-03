// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"bytes"
	"image"
	"io"

	"github.com/ronsor/majokko/format/png"
)

func init() {
	RegisterCodec(&PNGCodec{})
	RegisterCodec(&PNGCodec{isZNGCodec: true})
}

// PNGCodec is the still PNG codec.
type PNGCodec struct {
	isZNGCodec bool
}

// New returns a new instance of PNGCodec.
func (c *PNGCodec) New() Codec {
	return &PNGCodec{isZNGCodec: c.isZNGCodec}
}

// Name returns the name of the PNG codec: "png"
func (c *PNGCodec) Name() string {
	if c.isZNGCodec {
		return "zng"
	} else {
		return "png"
	}
}

// Magic returns magic strings that identify PNG data.
func (c *PNGCodec) Magic() []string {
	return []string{"\x89PNG\r\n\x1a\n"}
}

// textEntryToPNGChunk converts a TextEntry to a PNG
// `tEXt` or `iTXt` chunk.
// TODO: support compression
func textEntryToPNGChunk(te *TextEntry) (chunk png.Chunk) {
	var buf bytes.Buffer
	buf.WriteString(te.Key)
	buf.WriteByte(0)

	if te.Language == "" && te.Utf8Key == "" && !te.IsUtf8 {
		buf.WriteString(te.Value)

		chunk.Name = "tEXt"
	} else {
		buf.WriteByte(0)
		buf.WriteByte(0)

		buf.WriteString(te.Language)
		buf.WriteByte(0)

		buf.WriteString(te.Utf8Key)
		buf.WriteByte(0)

		buf.WriteString(te.Value)

		chunk.Name = "iTXt"
	}

	chunk.Data = buf.Bytes()
	return
}

// pngChunkToTextEntry converts a PNG `tEXt` or `iTXt`
// chunk to a TextEntry.
// TODO: support compression
func pngChunkToTextEntry(chunk png.Chunk) (te *TextEntry, err error) {
	te = &TextEntry{}

	keyNul := bytes.IndexByte(chunk.Data, 0)
	if keyNul == -1 {
		err = png.FormatError("invalid text-type chunk: " + chunk.Name)
		return
	} else if keyNul == (len(chunk.Data) - 1) {
		err = png.FormatError("truncated text-type chunk: " + chunk.Name)
		return
	}
	te.Key = string(chunk.Data[:keyNul])

	switch chunk.Name {
		case "tEXt":
			te.Value = string(chunk.Data[keyNul+1:])
			return
		case "iTXt":
			te.IsUtf8 = true

			rest := chunk.Data[keyNul+1:]
			if rest[0] != 0 {
				err = png.UnsupportedError("compressed iTXt chunks")
				return
			}

			if len(rest) < 4 {
				err = png.FormatError("truncated iTXt chunk")
				return
			}

			rest = rest[2:]
			nul := bytes.IndexByte(rest, 0)
			if nul == -1 {
				err = png.FormatError("truncated iTXt chunk: missing language")
				return
			} else if nul == (len(rest) - 1) {
				err = png.FormatError("truncated iTXt chunk after language")
				return
			}

			te.Language = string(rest[:nul])
			rest = rest[nul+1:]

			nul = bytes.IndexByte(rest, 0)
			if nul == -1 {
				err = png.FormatError("truncated iTXt chunk: missing utf-8 key")
				return
			} else if nul == (len(rest) - 1) {
				err = png.FormatError("truncated iTXt chunk after utf-8 key")
				return
			}
			te.Utf8Key = string(rest[:nul])

			rest = rest[nul+1:]
			te.Value = string(rest)
			return
	}

	err = png.UnsupportedError("text-type chunk: " + chunk.Name)
	return
}

// Decode decodes a PNG according to the options specified.
func (c *PNGCodec) Decode(r io.Reader, o *DecodeOptions) (image.Image, error) {
	if o == nil { o = DefaultDecodeOptions() }

	pngOpt := &png.DecodeOptions{
		ParseUnknownChunk: func (c png.Chunk) error {
			if o.Metadata == nil { return nil }

			if c.Name == "tEXt" || c.Name == "iTXt" {
				entry, err := pngChunkToTextEntry(c)
				if err != nil && o.Strict { return err }
				if err == nil {
					if entry.Key != "__COMMENT__" {
						o.Metadata.Text.Add(entry)
					} else {
						o.Metadata.Comments = append(o.Metadata.Comments, entry.Value)
					}
				}
			}
			return nil
		},
	}

	return png.DecodeWithOptions(r, pngOpt)
}

// DecodeConfig returns the color model and dimensions of a PNG image
// without decoding the image.
func (c *PNGCodec) DecodeConfig(r io.Reader, o *DecodeOptions) (image.Config, error) {
	return png.DecodeConfig(r)
}

// Encode encodes a PNG according to the options specified.
func (c *PNGCodec) Encode(w io.Writer, i image.Image, o *EncodeOptions) error {
	convertCompressionLevel := func (i int) png.CompressionLevel {
		if i == -1 {
			return png.DefaultCompression
		} else if i == 0 {
			return png.NoCompression
		} else if i < 50 {
			return png.BestSpeed
		} else if i < 75 {
			return png.DefaultCompression
		} else {
			return png.BestCompression
		}
	}

	if o == nil { o = DefaultEncodeOptions() }

	enc := &png.Encoder{
		CompressionLevel: convertCompressionLevel(o.CompressionLevel),
		UseZstd: c.isZNGCodec,
	}

	pngOpt := &png.EncodeOptions{}
	if o.Metadata != nil {
		for _, entry := range o.Metadata.Text {
			pngOpt.CustomChunks = append(pngOpt.CustomChunks, textEntryToPNGChunk(entry))
		}

		for _, comment := range o.Metadata.Comments {
			pngOpt.CustomChunks = append(pngOpt.CustomChunks, textEntryToPNGChunk(&TextEntry{
				Key: "__COMMENT__",
				Value: comment,
				IsUtf8: true,
			}))
		}
	}

	return enc.EncodeWithOptions(w, i, pngOpt)
}

var (
	_ Decoder = &PNGCodec{}
	_ Encoder = &PNGCodec{}
)
