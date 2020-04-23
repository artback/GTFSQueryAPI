package direction

import (
	"database/sql"
	"github.com/artback/gtfsQueryGoApi/query"
	"github.com/artback/gtfsQueryGoApi/time"
	"github.com/cornelk/hashmap"
	"log"
	"strconv"
	"time"
)

func GetResult(r *query.Repository, la float64, lo float64, radius int, maxDepartures int, maxStops int) []Result {
	rows, err := getStops(r.Db, strconv.FormatFloat(la, 'f', -1, 64),
		strconv.FormatFloat(lo, 'f', -1, 64), strconv.Itoa(radius))
	if err != nil {
		panic(err)
	}
	return groupAndSortRows(rows, maxStops, maxDepartures)
}

func getStops(db *sql.DB, lat string, lon string, radius string) (*sql.Rows, error) {
	return db.Query(
		"SELECT s.stop_id as id, arrival_time, departure_time, stop_name as name, stop_lat as lat, stop_lon as lon,trip_headsign as headsign, date" +
			" from stop_times JOIN stops s ON s.stop_id = stop_times.stop_id" +
			" JOIN trips t on stop_times.trip_id = t.trip_id JOIN calendar_dates cd on t.service_id = cd.service_id" +
			" WHERE(date(current_timestamp + interval'- 4 hours') = cd.date OR date(current_timestamp + interval '20 hours') = cd.date)" +
			" AND st_dwithin(geography(st_point(s.stop_lon, s.stop_lat)), geography(st_point( " + lon + " ," + lat + " ))," + radius + ") " +
			" ORDER BY cd.date, departure_time")
}

func groupAndSortRows(rows *sql.Rows, maxStops int, maxDepartures int) []Result {
	resultMap := hashmap.New(uintptr(maxStops * 2))
	for rows.Next() {
		var row row
		if err := rows.Scan(&row.id, &row.arrival_time, &row.departure_time, &row.name, &row.lat, &row.lon, &row.headsign, &row.date); err != nil {
			log.Fatal(err)
		}

		loc, _ := time.LoadLocation("Europe/Stockholm")
		timeDiff := time_processing.GetTimeDifference(loc, time.UTC)
		now := time.Now().In(time.UTC)
		date, _ := time.Parse(time.RFC3339, row.date)
		dep := time_processing.AddTime(date, row.departure_time).Add(time.Hour * time.Duration(-timeDiff))
		arr := time_processing.AddTime(date, row.arrival_time).Add(time.Hour * time.Duration(-timeDiff))
		if dep.After(now) {
			value, exist := resultMap.Get(row.id)
			if exist == true {
				v := value.(Result)
				if len(v.Departures) < maxDepartures {
					v.Departures = append(v.Departures, Departure{dep.Format("15:04:05"), arr.Format("15:04:05"), dep.Format("2006-01-02T15:04:05-07:00"), Trip{row.headsign}})
					resultMap.Set(row.id, v)
				}
			} else {
				if resultMap.Len() < maxStops {
					resultMap.Insert(row.id, Result{
						Stop{row.id, []string{row.lat, row.lon}, row.name},
						[]Departure{{dep.Format("15:04:05"), arr.Format("15:04:05"), dep.Format("2006-01-02T15:04:05-07:00"),
							Trip{row.headsign}}}})
				}
			}
		}
	}
	var r []Result
	for v := range resultMap.Iter() {
		r = append(r, v.Value.(Result))
	}
	return r
}
