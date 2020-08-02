package time_processing

import (
	"strconv"
	"strings"
	"time"
)

func AddTime(t time.Time, timeString string) (*time.Time, error) {
	parts := strings.Split(timeString, ":")
	var intParts []time.Duration
	for i := range parts {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, err
		}
		intParts = append(intParts, time.Duration(val))
	}
	t = t.Add(time.Hour * intParts[0]).Add(time.Minute * intParts[1]).Add(time.Second * intParts[2])
	return &t, nil
}

func GetTimeDifference(location *time.Location, location2 *time.Location) int {
	now := time.Now().In(location)
	return now.Hour() - time.Now().In(location2).Hour()
}
