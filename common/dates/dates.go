package dates

import (
	"time"
)

const (
	FmtDate              = "2006-01-02"
	FmtTime              = "15:04:05"
	FmtDateTime          = "2006-01-02 15:04:05"
	FmtDateTimeNoSeconds = "2006-01-02 15:04"
)

func NowDate() string {
	return Format(time.Now(), FmtDate)
}

func NowTime() string {
	return Format(time.Now(), FmtTime)
}

func NowDateTime() string {
	return Format(time.Now(), FmtDateTime)
}

func NowDateTimeNoSeconds() string {
	return Format(time.Now(), FmtDateTimeNoSeconds)
}

func NowUnix() int64 {
	return time.Now().Unix()
}

func NowTimestamp() int64 {
	return time.Now().UnixMilli()
}

func Format(time time.Time, layout string) string {
	return time.Format(layout)
}

func Parse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}
