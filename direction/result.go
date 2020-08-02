package direction

import (
	"database/sql"
	"fmt"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/allbin/gtfsQueryGoApi/time_processing"
	"github.com/cornelk/hashmap"
	"log"
	"time"
)

func GetResult(r *query.Repository, la float64, lo float64, radius int, maxDepartures int, maxStops int) ([]Result, error) {
	rows, err := r.GetStops(la, lo, radius, maxStops)
	if err != nil {
		return nil, err
	}
	return groupSort(rows, maxStops, maxDepartures)
}

func groupSort(rows *sql.Rows, maxStops int, maxDepartures int) ([]Result, error) {
	resultMap := hashmap.New(uintptr(maxStops * 2))

	for rows.Next() {
		var row row
		if err := rows.Scan(&row.id, &row.arrivalTime, &row.departureTime, &row.name, &row.lat, &row.lon, &row.headsign, &row.date, &row.dateString); err != nil {
			log.Fatal(err)
		}
		times, err := getTimes(row)
		if err != nil {
			return nil, err
		}

		now := time.Now().In(time.UTC)
		if times.Departure.After(now) {
			value, exist := resultMap.Get(row.id)
			if exist == true {
				v := value.(Result)
				if len(v.Departures) < maxDepartures {
					v.Departures = append(v.Departures, newDeparture(row, times))
					resultMap.Set(row.id, v)
				}

			} else {
				resultMap.Insert(row.id, rowResult(row, times))

			}
		}
	}
	var r []Result
	for v := range resultMap.Iter() {
		r = append(r, v.Value.(Result))
	}
	return r, nil
}
func getTimes(row row) (*times, error) {
	locName := "Europe/Stockholm"
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return nil, fmt.Errorf("problem loading location %s", locName)
	}
	timeDiff := time_processing.GetTimeDifference(loc, time.UTC)
	date, _ := time.Parse(time.RFC3339, row.date)
	departureTime, err := time_processing.AddTime(date, row.departureTime)
	if err != nil {
		return nil, err
	}
	dep := departureTime.Add(time.Hour * time.Duration(-timeDiff))
	arrivalTime, err := time_processing.AddTime(date, row.arrivalTime)
	if err != nil {
		return nil, err
	}
	arr := arrivalTime.Add(time.Hour * time.Duration(-timeDiff))
	return &times{arr, dep}, nil
}
func rowResult(row row, times *times) Result {
	return Result{
		Stop{
			row.id,
			[]string{row.lat, row.lon},
			row.name,
		},
		[]Departure{
			{
				times.Departure.Format("15:04:05"),
				times.Arrival.Format("15:04:05"),
				times.Departure.Format("2006-01-02T15:04:05-07:00"),
				Trip{row.headsign},
			},
		},
	}
}
func newDeparture(row row, times *times) Departure {
	return Departure{
		times.Departure.Format("15:04:05"),
		times.Arrival.Format("15:04:05"),
		times.Departure.Format("2006-01-02T15:04:05-07:00"),
		Trip{row.headsign},
	}
}
