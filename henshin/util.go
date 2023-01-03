// Copyright 2023 Ronsor Labs & The Go Authors. All rights reserved.

package henshin

// fmtExpand replaces %[var] or %v in the string based on the mapping function.
func fmtExpand(s string, mapping func(string) string) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '%' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := getFmtName(s[j+1:])
			if name == "" && w > 0 {
				// Encountered invalid syntax; eat the
				// characters.
			} else if name == "" {
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			} else {
				buf = append(buf, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

func getFmtName(s string) (string, int) {
	switch {
	case s[0] == '[':
		for i := 1; i < len(s); i++ {
			if s[i] == ']' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "%[]"
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "%["
	default:
		return s[0:1], 1
	}
}
