package strcase

import (
	"strings"
)


func ToSnake(s string) string {
	return ToDelimited(s, '_')
}


func ToScreamingSnake(s string) string {
	return ToScreamingDelimited(s, '_', true)
}


func ToKebab(s string) string {
	return ToDelimited(s, '-')
}

func ToScreamingKebab(s string) string {
	return ToScreamingDelimited(s, '-', true)
}

func ToDelimited(s string, del uint8) string {
	return ToScreamingDelimited(s, del, false)
}

func ToScreamingDelimited(s string, del uint8, screaming bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {

		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			if (v >= 'A' && v <= 'Z' && next >= 'a' && next <= 'z') || (v >= 'a' && v <= 'z' && next >= 'A' && next <= 'Z') {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != del && nextCaseIsChanged {

			if v >= 'A' && v <= 'Z' {
				n += string(del) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(del)
			}
		} else if v == ' ' || v == '_' || v == '-' {

			n += string(del)
		} else {
			n = n + string(v)
		}
	}

	if screaming {
		n = strings.ToUpper(n)
	} else {
		n = strings.ToLower(n)
	}
	return n
}
