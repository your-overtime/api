package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SafeGetUInt(i interface{}) uint {
	switch i.(type) {
	case uint32:
		return uint(i.(uint32))
	case uint64:
		return uint(i.(uint64))
	case int32:
		return uint(i.(int32))
	case int64:
		return uint(i.(int64))
	case float32:
		return uint(i.(float32))
	case float64:
		return uint(i.(float64))
	default:
		return uint(0)
	}
}

func DayStart(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
}

func DayEnd(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day()+1, 0, 0, 0, 0, day.Location())
}
