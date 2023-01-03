// Copyright 2022 Ronsor Labs. All rights reserved.

package henshin

type TextEntry struct {
	Key, Value string

	Language, Utf8Key string
	IsUtf8, Compress bool
}

type TextData []*TextEntry

func NewTextData() TextData {
	return make(TextData, 0, 8)
}

func NewTextDataFromStringMap(strMap map[string]string) TextData {
	td := make(TextData, 0, len(strMap))
	for k, v := range strMap {
		td.AddString(k, v)
	}
	return td
}

func (td *TextData) AddString(key, value string) {
	td.Add(&TextEntry{Key: key, Value: value})
}

func (td *TextData) Add(te *TextEntry) {
	*td = append(*td, te)
}

func (td TextData) GetString(key string) (value string, ok bool) {
	teValue, ok := td.Get(key)
	if ok {
		value = teValue.Value
	}
	return
}

func (td TextData) Get(key string) (value *TextEntry, ok bool) {
	for _, v := range td {
		if v.Key == key {
			value = v
			ok = true
			return
		}
	}
	return
}

func (td *TextData) Set(te *TextEntry) (replaced bool) {
	for i, v := range *td {
		if v.Key == te.Key {
			(*td)[i] = te
			replaced = true
			return
		}
	}

	td.Add(te)
	return
}

func (td *TextData) SetString(key, value string) (replaced bool) {
	return td.Set(&TextEntry{Key: key, Value: value})
}

func (td TextData) ToStringMap() (ret map[string]string) {
	ret = map[string]string{}
	for _, v := range td {
		ret[v.Key] = v.Value
	}
	return
}

// Metadata is additional image metadata.
type Metadata struct {
	// Text contains key-value pairs of text data.
	Text TextData

	// Comments contains a list of comments for the image.
	Comments []string

	// Specific contains encoder/decoder specific data.
	Specific any
}

func (md *Metadata) Clone() *Metadata {
	text := TextData(nil)
	for _, v := range md.Text {
		text = append(text, v)
	}
	comments := []string(nil)
	for _, v := range md.Comments {
		comments = append(comments, v)
	}

	return &Metadata{
		Text: text,
		Comments: comments,
		Specific: md.Specific,
	}
}
