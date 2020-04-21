package direction

import (
	"database/sql"
	"github.com/artback/gtfsQueryGoApi/query"
	"github.com/cornelk/hashmap"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetResult(r *query.Repository, la float64, lo float64, radius int, maxDepartures int, maxStops int) []Result {
	rows, err := getstops(r.Db, strconv.FormatFloat(la, 'f', -1, 64),
		strconv.FormatFloat(lo, 'f', -1, 64), strconv.Itoa(radius))
	if err != nil {
		panic(err)
	}
	return groupAndSortRows(rows, maxStops, maxDepartures)
}

func getstops(db *sql.DB, lat string, lon string, radius string) (*sql.Rows, error) {
	return db.Query(
		"SELECT s.stop_id as id, arrival_time, departure_time, stop_name as name, stop_lat as lat, stop_lon as lon,trip_headsign as headsign, date" +
			" from stop_times JOIN stops s ON s.stop_id = stop_times.stop_id" +
			" JOIN trips t on stop_times.trip_id = t.trip_id JOIN calendar_dates cd on t.service_id = cd.service_id" +
			" WHERE(date(current_timestamp + interval'- 4 hours') = cd.date OR date(current_timestamp + interval '20 hours') = cd.date)" +
			" AND st_dwithin(geography(st_point(s.stop_lon, s.stop_lat)), geography(st_point( " + lon + " ," + lat + " ))," + radius + ") " +
			" ORDER BY cd.date, departure_time")
}

func groupAndSortRows(rows *sql.Rows, maxStops int, maxDepartures int) []Result {
	hashMap := hashmap.New(uintptr(maxStops * 2))
	for rows.Next() {
		var row row
		if err := rows.Scan(&row.id, &row.arrival_time, &row.departure_time, &row.name, &row.lat, &row.lon, &row.headsign, &row.date); err != nil {
			log.Fatal(err)
		}
		loc, _ := time.LoadLocation("Europe/Stockholm")
		now := time.Now().In(loc)
		hour_diff := now.Hour() - time.Now().In(time.UTC).Hour()
		d, _ := time.Parse(time.RFC3339, row.date)
		d = addTime(d, row.departure_time).Add(time.Hour * time.Duration(-hour_diff))
		if d.After(now) {
			if hashMap.Len() < maxStops {
				value, exist := hashMap.GetOrInsert(row.id,
					Result{
						Stop{row.id, []string{row.lat, row.lon}, row.name},
						[]Departure{{row.departure_time, row.arrival_time, d.Format(time.RFC3339),
							Trip{row.headsign}}}})
				if exist == true {
					v := value.(Result)
					if len(v.Departures) < maxDepartures {
						v.Departures = append(v.Departures, Departure{row.departure_time, row.arrival_time, d.Format(time.RFC3339), Trip{row.headsign}})
						hashMap.Set(row.id, v)
					}
				}
			}
		}
	}
	r := []Result{}
	for v := range hashMap.Iter() {
		r = append(r, v.Value.(Result))
	}
	return r
}

func addTime(t time.Time, tstring string) time.Time {
	parts := strings.Split(tstring, ":")
	int_parts := []time.Duration{}
	for i, _ := range parts {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			panic(err)
		}
		int_parts = append(int_parts, time.Duration(val))
	}
	return t.Add(time.Hour * int_parts[0]).Add(time.Minute * int_parts[1]).Add(time.Second * int_parts[2])
}
