package direction

type result struct {
	stop       stop
	departures []departures
}

type stop struct {
	id   string
	loc  []float64
	name string
}
type departures struct {
	departure_time string
	arrival_time   string
	date           string
	trip           trip
}
type trip struct {
	headsign string
}
