package util

import (
	"fmt"
	"time"
)

func Now() time.Time {
	return time.Now().UTC()
}

func NowUnix() int64 {
	return Now().Unix()
}

func TimeFromUnix(epoch int64) time.Time {
	return time.Unix(epoch, 0).UTC()
}

func TimeDiffDays(tick, tok time.Time) int {
	return int(tick.Sub(tok).Hours() / 24)
}

func DaysFromNowToTime(tok time.Time) int {
	return TimeDiffDays(Now(), tok)
}

func DaysFromNowToTimeStamp(tok int64) int {
	return TimeDiffDays(TimeFromUnix(tok), Now())
}

func HumanReadableDate(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("Monday, 02 January 2006 15:04:05 MST")
}

// ParseTimeString parses an RFC3339 timestamp string into a time.Time object
func ParseTimeString(timestamp string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp format: %w", err)
	}
	return parsedTime, nil
}
