package utils

import (
	"time"
)

func TimeFormat() string {
	return "2006-01-02 15:04:05"
}

func GetNowTimeCST() (time.Time, string) {
	NowTimeZone := time.FixedZone("CST", 8*3600)
	nowAt := time.Now().In(NowTimeZone)
	fmtNow := nowAt.Format(TimeFormat())
	return nowAt, fmtNow
}

func GetMinTime() (time.Time,string) {
	return time.Unix(0,0), time.Unix(0,0).Format(TimeFormat())
}