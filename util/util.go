package util

import (
	"fmt"
	"time"
)

func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func ShortDate(date *time.Time) string {
	now := time.Now()
	ret := fmt.Sprintf("%02d시 %02d분", date.Hour(), date.Minute())
	need := 0
	if now.Year() != date.Year() {
		need = 3
	} else if now.Month() != date.Month() {
		need = 2
	} else if now.Day() != date.Day() {
		need = 1
	}
	if need >= 1 {
		ret = fmt.Sprintf("%d일 %s", now.Day(), ret)
	}
	if need >= 2 {
		ret = fmt.Sprintf("%d월 %s", now.Month(), ret)
	}
	if need >= 3 {
		ret = fmt.Sprintf("%d년 %s", now.Year(), ret)
	}
	return ret
}
