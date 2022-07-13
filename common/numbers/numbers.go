package numbers

import "strconv"

func ToInt64(str string) int64 {
	return ToInt64ByDefault(str, 0)
}

func ToInt64ByDefault(str string, def int64) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		val = def
	}
	return val
}

func ToInt(str string) int {
	return ToIntByDefault(str, 0)
}

func ToIntByDefault(str string, def int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = def
	}
	return val
}
