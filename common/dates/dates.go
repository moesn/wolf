package dates

import (
	"strconv"
	"time"
)

const (
	FmtDate              = "2006-01-02"
	FmtTime              = "15:04:05"
	FmtDateTime          = "2006-01-02 15:04:05"
	FmtDateTimeNoSeconds = "2006-01-02 15:04"
)


func NowUnix() int64 {
	return time.Now().Unix()
}

func NowDateTime() string {
	return Format(time.Now(),FmtDateTime)
}

func FromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func NowTimestamp() int64 {
	return Timestamp(time.Now())
}

func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func FromTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

func Format(time time.Time, layout string) string {
	return time.Format(layout)
}

func Parse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

func GetDay(time time.Time) int {
	ret, _ := strconv.Atoi(time.Format("20060102"))
	return ret
}

func WithTimeAsStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
