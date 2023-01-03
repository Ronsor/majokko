// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"bytes"
	"image"

	"testing"
)

// TestTextChunksRoundTrip tests the encoding and decoding of
// TextEntry objects to and from PNG chunks.
func TestTextChunksRoundTrip(t *testing.T) {
	textEntryEq := func (a, b *TextEntry) bool {
		return a.Key == b.Key && a.Value == b.Value &&
			a.Compress == b.Compress &&
			a.Utf8Key == b.Utf8Key && a.Language == b.Language &&
			(a.IsUtf8 == b.IsUtf8 ||
			b.IsUtf8 && a.Language != "" ||
			b.IsUtf8 && a.Utf8Key != "")
	}

	cases := []*TextEntry{
		&TextEntry{Key: "simple", Value: "this is some text"},
		&TextEntry{Key: "spaces and NUL", Value: "this is some text\x00with an embedded NUL"},
		&TextEntry{Key: "utf8 stuff", Value: "things here\x00 and a NUL", IsUtf8: true},
		&TextEntry{
			Key: "kitchen sink",
			Utf8Key: "utf8-y kitchen sink",
			Value: "\x00here",
			Language: "en_US",
		},
		&TextEntry{
			Key: "",
			Utf8Key: "huh",
			Value: "test123",
			
		},
	}

	for _, original := range cases {
		chunk := textEntryToPNGChunk(original)
		decoded, err := pngChunkToTextEntry(chunk)
		if err != nil {
			t.Fatal(err)
		}

		if !textEntryEq(original, decoded) {
			t.Logf(`raw data: %v`, chunk.Data)
			t.Fatalf(`original and decoded TextEntry do not match: %v != %v`, original, decoded)
		}
	}
}

// TestFullRoundTrip tests the encoding and decoding of
// PNG images with text metadata.
func TestFullRoundTrip(t *testing.T) {
	img := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{512, 512},
	})

	textEntryEq := func (a, b *TextEntry) bool {
		return a.Key == b.Key && a.Value == b.Value &&
			a.Compress == b.Compress &&
			a.Utf8Key == b.Utf8Key && a.Language == b.Language &&
			(a.IsUtf8 == b.IsUtf8 ||
			b.IsUtf8 && a.Language != "" ||
			b.IsUtf8 && a.Utf8Key != "")
	}

	text := TextData{
		&TextEntry{Key: "simple", Value: "this is some text"},
		&TextEntry{Key: "spaces and NUL", Value: "this is some text\x00with an embedded NUL"},
		&TextEntry{Key: "utf8 stuff", Value: "things here\x00 and a NUL", IsUtf8: true},
		&TextEntry{
			Key: "kitchen sink",
			Utf8Key: "utf8-y kitchen sink",
			Value: "\x00here",
			Language: "en_US",
		},
		&TextEntry{
			Key: "",
			Utf8Key: "huh",
			Value: "test123",
			
		},
	}

	comments := []string{
		"first comment",
		"second one - something else here",
		"a bunch of stuff?",
		"last one",
	}

	encOpt := &EncodeOptions{
		Metadata: &Metadata{
			Text: text,
			Comments: comments,
		},
	}

	var codec PNGCodec

	var buf bytes.Buffer
	err := codec.Encode(&buf, img, encOpt)
	if err != nil { t.Fatal(err) }

	decOpt := &DecodeOptions{
		Metadata: &Metadata{},
	}
	decodedImg, err := codec.Decode(&buf, decOpt)

	for i, v := range decOpt.Metadata.Text {
		if !textEntryEq(text[i], v) {
			t.Fatalf(`original and decoded TextEntry do not match: %v != %v`, text[i], v)
		}
	}

	for i, v := range decOpt.Metadata.Comments {
		if comments[i] != v {
			t.Fatalf(`original and decoded comments do not match: %q != %q`, comments[i], v)
		}
	}

	decodedImgNRGBA := decodedImg.(*image.NRGBA)
	for i, v := range decodedImgNRGBA.Pix {
		if img.Pix[i] != v {
			t.Fatalf(`original and decoded image do not match at %d: %v != %v`, i, img.Pix[i], v)
		}
	}
}
