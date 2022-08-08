package strs

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/moesn/wolf/common/jsons"
	"github.com/moesn/wolf/common/structs"
	"strconv"
	"strings"
	"unicode"
)

func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func IsAnyBlank(strs ...string) bool {
	for _, str := range strs {
		if IsBlank(str) {
			return true
		}
	}
	return false
}

func DefaultIfBlank(str, def string) string {
	if IsBlank(str) {
		return def
	} else {
		return str
	}
}

func IsEmpty(str string) bool {
	return len(str) == 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

func Equals(a, b string) bool {
	return a == b
}

func EqualsIgnoreCase(a, b string) bool {
	return a == b || strings.ToUpper(a) == strings.ToUpper(b)
}

func UUID() string {
	u, _ := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

func ToString(i interface{}) string {
	str:=""
	switch i.(type) {
	case string:
		str= i.(string)
	case bool:
		str= strconv.FormatBool(i.(bool))
	case int:
		str= strconv.Itoa(i.(int))
	case int64:
		str= strconv.FormatInt(i.(int64),10)
	case float64:
		str= strconv.FormatFloat(i.(float64),'f', 10, 64)
	case structs.JSON:
		str= jsons.ToString(i)
	}

	return str
}
