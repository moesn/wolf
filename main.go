package main

import (
	"fmt"
	"github.com/moesn/wolf/common/dates"
)

func main() {
	fmt.Println(dates.NowDateTimeNoSeconds())
	fmt.Println(dates.NowDateTime())
	fmt.Println(dates.NowDate())
	fmt.Println(dates.NowTime())
	fmt.Println(dates.NowUnix())
	fmt.Println(dates.NowTimestamp())
	fmt.Println(dates.Parse("2015-01-01 12:12:12",dates.FmtDateTime))
}
