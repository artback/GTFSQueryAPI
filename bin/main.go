package main

import (
	"github.com/artback/gtfsQueryGoApi/config"
	"github.com/artback/gtfsQueryGoApi/direction"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	conf *config.Configuration
)

func init() {
	conf = new(config.Configuration)
	err := config.Init(conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/departures/place", direction.PlaceHandler).Queries("lat", "{lat}",
		"lon", "{lon}", "adress", "{adress}", "radius", "{radius}", "maxDepartures", "{maxDepartures}", "maxStops", "{maxStops}")
	log.Fatal(http.ListenAndServe(":8080", r))
}
