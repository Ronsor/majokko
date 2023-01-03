// Copyright 2023 Ronsor Labs. All rights reserved.

package henshin

import (
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

var emptyImage = image.NewRGBA(image.Rect(0, 0, 0, 0))

type Wand struct {
	im image.Image
	md *Metadata

	decOpt *DecodeOptions
	encOpt *EncodeOptions
}

func NewWand() *Wand {
	md := &Metadata{}

	return &Wand{
		im: nil,
		md: md,
		decOpt: &DecodeOptions{
			Metadata: md,
		},
		encOpt: &EncodeOptions{
			CompressionLevel: -1,
			Metadata: md,
		},
	}
}

func (w *Wand) NewImage(iw, ih int) {
	w.im = image.NewRGBA(image.Rect(0, 0, iw, ih))
}

func (w *Wand) SetImage(im image.Image) {
	w.im = im
}

func (w *Wand) Image() image.Image {
	return w.im
}

func (w *Wand) DecodeImage(r io.Reader) error {
	im, err := Decode(r, w.decOpt)
	if err != nil { return err }
	w.im = im
	return nil
}

func (w *Wand) EncodeImage(wr io.Writer, codec string) error {
	if w.im == nil {
		return Encode(codec, wr, emptyImage, w.encOpt)
	}
	return Encode(codec, wr, w.im, w.encOpt) 
}

func (w *Wand) ReadImage(path string) error {
	if path == "-" {
		return w.DecodeImage(os.Stdin)
	}

	f, err := os.Open(path)
	if err != nil { return err }
	return w.DecodeImage(f)
}

func (w *Wand) WriteImage(path string) error {
	codecName := "png"

	hasColon := strings.IndexByte(path, ':')
	hasExt := strings.LastIndexByte(path, '.')
	if hasColon != -1 {
		_, err := NewCodec(path[:hasColon])
		if err == nil {
			codecName = path[:hasColon]
			path = path[hasColon+1:]
		}
	} else if hasExt != -1 {
		ext := path[hasExt+1:]
		_, err := NewCodec(ext)
		if err == nil {
			codecName = ext
		}
	}

	if path == "-" {
		return w.EncodeImage(os.Stdout, codecName)
	}

	f, err := os.Create(path)
	if err != nil { return err }
	return w.EncodeImage(f, codecName)
}

func (w *Wand) AddComment(comment string) {
	w.md.Comments = append(w.md.Comments, comment)
}

func (w *Wand) SetComments(comments []string) {
	w.md.Comments = comments
}

func (w *Wand) Comments() []string {
	return w.md.Comments
}

func (w *Wand) Metadata() *Metadata {
	return w.md
}

func (w *Wand) Strip() {
	*w.md = Metadata{}
}

func (w *Wand) Width() int {
	if w.im == nil { return 0 }
	return w.im.Bounds().Dx()
}

func (w *Wand) Height() int {
	if w.im == nil { return 0 }
	return w.im.Bounds().Dy()
}

func (w *Wand) Resize(iw, ih int, strategy ResizeStrategy) {
	if w.Width() == iw && w.Height() == ih { return }

	newIm := image.NewRGBA(image.Rect(0, 0, iw, ih))
	if (w.Width() == 0 && w.Height() == 0) || (iw == 0 && ih == 0) {
		w.im = newIm
		return
	}

	strategy.Scale(newIm, newIm.Bounds(), w.im, w.im.Bounds(), draw.Over, nil)
	w.im = newIm
}

func (w *Wand) ResizeMaxArea(area int, strategy ResizeStrategy) {
	iw, ih := areaFit(w.Width(), w.Height(), area)
	w.Resize(iw, ih, strategy)
}

func (w *Wand) Crop(iw, ih, xoff, yoff int) {
	if w.Width() == iw && w.Height() == ih && xoff == 0 && yoff == 0 { return }

	newIm := image.NewRGBA(image.Rect(0, 0, iw, ih))
	if (w.Width() == 0 && w.Height() == 0) || (iw == 0 && ih == 0) {
		w.im = newIm
		return
	}

	draw.Copy(newIm, newIm.Bounds().Min, w.im, image.Rect(xoff, yoff, xoff + iw, yoff + ih), draw.Over, nil)
	w.im = newIm
}

func (w *Wand) CropAnchor(iw, ih int, anchor string) {
	panic("TODO")
}

func (w *Wand) ForceRGBA() {
	if w.im != nil {
		newIm := image.NewRGBA(w.im.Bounds())
		draw.Copy(newIm, newIm.Bounds().Min, w.im, w.im.Bounds(), draw.Over, nil)
		w.im = newIm
	}
}

func (w *Wand) Clone() *Wand {
	var newIm draw.Image

	if w.im != nil {
		newIm = image.NewRGBA(w.im.Bounds())
		draw.Copy(newIm, newIm.Bounds().Min, w.im, w.im.Bounds(), draw.Over, nil)
	}

	newMd := w.md.Clone()

	return &Wand{
		im: newIm,
		md: newMd,

		decOpt: &DecodeOptions{
			Metadata: newMd,
		},
		encOpt: &EncodeOptions{
			Metadata: newMd,
			CompressionLevel: w.encOpt.CompressionLevel,
		},
	}
}

func (w *Wand) Hash() uint64 {
	if w.im == nil || (w.Width() == 0 && w.Height() == 0) { return 0 }

	smallIm := image.NewGray(image.Rect(0, 0, 9, 8))
	NearestStrategy.Scale(smallIm, smallIm.Bounds(), w.im, w.im.Bounds(), draw.Over, nil)
	return diffHash(smallIm)
}

func (w *Wand) SetCompressionQuality(q int) {
	w.encOpt.CompressionLevel = 100 - q
}

func (w *Wand) SetCompressionLevel(l int) {
	w.encOpt.CompressionLevel = l
}

func (w *Wand) property(key string) (val string) {
	switch (key) {
		case "w", "width": val = strconv.Itoa(w.Width())
		case "h", "height": val = strconv.Itoa(w.Height())
		case "H", "hash": val = strconv.FormatInt(int64(w.Hash()), 10)
		case "J", "json":
			val = `{"width":` + strconv.Itoa(w.Width()) + `,"height":` + strconv.Itoa(w.Height()) + `,"hash":` + strconv.FormatInt(int64(w.Hash()), 10) + `}`
	}
	return
}

func (w *Wand) FormatString(fmt string) (ret string) {
	return fmtExpand(fmt, w.property)
}
