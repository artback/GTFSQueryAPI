package query

import (
	"database/sql"
	"fmt"
	"github.com/allbin/gtfsQueryGoApi/internal/config"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type Repository struct{ Db *sql.DB }

func (r *Repository) Connect(c config.DatabaseConfiguration) error {
	passwordArg := ""
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = c.Password
	}
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = c.Host
	}
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = c.Port
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = c.User
	}
	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = c.Database
	}
	if len(pass) > 0 {
		passwordArg = "password=" + pass
	}
	var err error
	dbString := fmt.Sprintf("host=%s port=%d user=%s %s dbname=%s sslmode=disable",
		host, port, user, passwordArg, db)
	r.Db, err = sql.Open(c.Driver, dbString)
	if err != nil {
		return err
	}

	return r.Db.Ping()
}

func NewConnected(c config.DatabaseConfiguration) (*Repository, error) {
	repo := new(Repository)
	err := repo.Connect(c)
	return repo, err
}

func (r *Repository) GetStops(la float64, lo float64, radius int, maxStops int) (*sql.Rows, error) {
	return r.Db.Query(stopsQuery(la, lo, radius, maxStops))
}
func stopsQuery(la float64, lo float64, radius int, maxStops int) string {
	lat := strconv.FormatFloat(la, 'f', -1, 64)
	lon := strconv.FormatFloat(lo, 'f', -1, 64)
	r := strconv.Itoa(radius)
	m := strconv.Itoa(maxStops)
	return fmt.Sprintf(`SELECT s.stop_id as id, arrival_time, departure_time, stop_name as name, stop_lat as lat, stop_lon as lon,
			 trip_headsign as headsign, date, (date::varchar || ' ' || departure_time) as date_string
			 from stop_times JOIN stops s ON s.stop_id = stop_times.stop_id
			 JOIN trips t on stop_times.trip_id = t.trip_id JOIN calendar_dates cd on t.service_id = cd.service_id
			 WHERE s.stop_id in (select distinct stop_id from stops where st_dwithin(geography(st_point(stop_lat, stop_lon)), geography(st_point(%s,%s)), %s)
			 order by stop_id limit %s)  AND ((date(current_timestamp + interval '- 4 hours') = cd.date
			 OR date(current_timestamp + interval '20 hours') = cd.date)) order by date_string;`, lat, lon, r, m)
}
