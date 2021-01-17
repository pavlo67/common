package common

import "time"

func RoundToDays(tm *time.Time) *time.Time {
	if tm == nil {
		return nil
	}

	tmUTC := tm.UTC()

	tmRounded := time.Date(tmUTC.Year(), tmUTC.Month(), tmUTC.Day(), 0, 0, 0, 0, time.UTC)

	return &tmRounded
}

func RoundToMilliseconds(tm *time.Time) *time.Time {
	if tm == nil {
		return nil
	}

	tmUTC := tm.UTC()
	microseconds := time.Duration(tmUTC.Nanosecond()) / (time.Nanosecond * 1000000)

	tmRounded := time.Date(tmUTC.Year(), tmUTC.Month(), tmUTC.Day(), tmUTC.Hour(), tmUTC.Minute(), tmUTC.Second(), int(microseconds), time.UTC)

	return &tmRounded
}
