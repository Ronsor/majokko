// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

// EncodeOptions specifies options for image encoders.
type EncodeOptions struct {
	// CompressionLevel is the compression level to use, on a scale
	// of 0 to 100. A value of 0 indicates no compression for PNG,
	// and 100% quality for JPEG. Other codecs may interpret this
	// value differently.
	CompressionLevel int

	// Metadata is any metadata or additional properties associated
	// with the image.
	Metadata *Metadata

	// EncoderSpecific is encoder-specific options.
	EncoderSpecific any
}

// DefaultEncodeOptions returns the default encoding options.
func DefaultEncodeOptions() *EncodeOptions {
	return &EncodeOptions{
		CompressionLevel: -1,
		Metadata: nil,
		EncoderSpecific: nil,
	}
}

// DecodeOptions specifies options for image decoders.
type DecodeOptions struct {
	// Metadata is any metadata or additional properties associated
	// with the image. Set this to an empty metadata struct in order
	// to read this data.
	Metadata *Metadata

	// Strict tells the decoder whether or not to strictly parse even
	// non-critical metadata.
	Strict bool

	// DecoderSpecific is decoder-specific options.
	DecoderSpecific any
}

// DefaultDecodeOptions returns the default decoding options.
func DefaultDecodeOptions() *DecodeOptions {
	return &DecodeOptions{}
}
