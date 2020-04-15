package direction

type Result struct {
	stop       Stop
	departures Departures
}

type Stop struct {
	id   string
	loc  []float64
	name string
}
type Departures struct {
	departure_time
}
