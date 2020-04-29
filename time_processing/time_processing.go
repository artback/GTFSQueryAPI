package time_processing

import (
	"strconv"
	"strings"
	"time"
)

func AddTime(t time.Time, tstring string) time.Time {
	parts := strings.Split(tstring, ":")
	var intParts []time.Duration
	for i := range parts {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			panic(err)
		}
		intParts = append(intParts, time.Duration(val))
	}
	return t.Add(time.Hour * intParts[0]).Add(time.Minute * intParts[1]).Add(time.Second * intParts[2])
}

func GetTimeDifference(location *time.Location, location2 *time.Location) int {
	now := time.Now().In(location)
	return now.Hour() - time.Now().In(location2).Hour()
}
