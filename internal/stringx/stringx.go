package stringx

import (
	"bytes"
	"unicode"
)

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func SnakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			_, _ = buf.WriteRune(c - 'A' + 'a')
		} else {
			_, _ = buf.WriteRune(c)
		}
	}
	return buf.String()
}
