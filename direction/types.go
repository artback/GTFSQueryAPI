package direction

type Result struct {
	Stop       Stop        `json:"stop"`
	Departures []Departure `json:"departures"`
}

type Stop struct {
	Id   string   `json:"id"`
	Loc  []string `json:"loc"`
	Name string   `json:"name"`
}
type Departure struct {
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Date          string `json:"date"`
	Trip          Trip   `json:"trip"`
}
type Trip struct {
	Headsign string `json:"headsign"`
}

type row struct {
	id             string
	arrival_time   string
	departure_time string
	name           string
	lat            string
	lon            string
	headsign       string
	date           string
}
