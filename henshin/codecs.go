// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

import (
	"bufio"
	"image"
	"io"
	"sort"
	"strings"
)

type ErrNoSuchCodec string
func (e ErrNoSuchCodec) Error() string { return "no such codec: " + string(e) }

var (
	knownCodecs = map[string]Codec{}
	knownCodecAliases = map[string]string{}
)

// Codec is any named image codec.
type Codec interface {
	Name() string
	New() Codec
}

// CodecWithAliases is an image codec that can have alternate
// names. Use this if you want to support format autodetection
// for multiple file extensions.
type CodecWithAliases interface {
	Codec
	Aliases() []string
}

// Decoder is an image codec that can decode.
type Decoder interface {
	Codec
	Magic() []string
	Decode(io.Reader, *DecodeOptions) (image.Image, error)
	DecodeConfig(io.Reader, *DecodeOptions) (image.Config, error)
}

// Encoder is an image codec that can encode.
type Encoder interface {
	Codec
	Encode(io.Writer, image.Image, *EncodeOptions) error
}

// RegisterCodec registers a new image processing codec.
func RegisterCodec(c Codec) {
	knownCodecs[c.Name()] = c
	if withAliases, yes := c.(CodecWithAliases); yes {
		for _, alias := range withAliases.Aliases() {
			knownCodecAliases[alias] = c.Name()
		}
	}
}

// NewCodec creates a new instance of the codec corresponding
// to the specified name.
func NewCodec(name string) (Codec, error) {
	alias, ok := knownCodecAliases[name]
	if ok { name = alias }
	codec, ok := knownCodecs[name]
	if !ok { return nil, ErrNoSuchCodec(name) }
	return codec.New(), nil
}

// Codecs returns a list of all codecs.
func Codecs() (ret []Codec) {
	for _, codec := range knownCodecs {
		ret = append(ret, codec)
	}
	sort.Slice(ret, func (i, j int) bool {
		return strings.Compare(ret[i].Name(), ret[j].Name()) < 0
	})
	return
}

// peekableReader is an io.Reader that implements the Peek()
// function.
type peekableReader interface {
	io.Reader
	Peek(n int) ([]byte, error)
}

func match(a []byte, b string) bool {
	if len(a) != len(b) { return false }

	for i, v := range a {
		if v == b[i] { continue }
		if b[i] == '?' { continue }
		return false
	}
	return true
}

// Decode decodes an image that has been encoded in a format understood
// by a registered codec.
func Decode(r io.Reader, o *DecodeOptions) (image.Image, error) {
	pkr, ok := r.(peekableReader)
	if !ok {
		pkr = bufio.NewReader(r)
	}

	for _, v := range knownCodecs {
		d, ok := v.(Decoder)
		if !ok { continue }

		for _, m := range d.Magic() {
			toPeek, err := pkr.Peek(len(m))
			if err == nil && match(toPeek, m) {
				return d.New().(Decoder).Decode(pkr, o)
			}
		}
	}
	return nil, ErrNoSuchCodec("unknown")
}

// Encode encodes an image into the registered codec corresponding to
// the specified name.
func Encode(name string, w io.Writer, i image.Image, o *EncodeOptions) error {
	codec, err := NewCodec(name)
	if err != nil { return err }
	encoder, ok := codec.(Encoder)
	if !ok { return ErrNoSuchCodec(name) }

	return encoder.Encode(w, i, o)
}
