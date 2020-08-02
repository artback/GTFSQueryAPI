package direction

import (
	"encoding/json"
	"fmt"
	"github.com/allbin/gtfsQueryGoApi/geocode"
	"github.com/allbin/gtfsQueryGoApi/internal/config"
	"github.com/allbin/gtfsQueryGoApi/query"
	geolocate "github.com/martinlindhe/google-geolocate"
	"net/http"
	"os"
	"strconv"
)

type Places struct {
	Repo *query.Repository
	config.DefaultConfiguration
	Geo *geolocate.GoogleGeo
}

func (p *Places) PlaceHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	apiKey := os.Getenv("GTFS_QUERY_API_KEY")
	if apiKey != "" {
		k := q.Get("k")
		if k != apiKey {
			http.Error(w, "Missing k parameter(API KEY)", http.StatusUnauthorized)
			return
		}

	}
	lat, _ := strconv.ParseFloat(q.Get("lat"), 32)
	lon, _ := strconv.ParseFloat(q.Get("lon"), 32)

	radius, _ := strconv.ParseInt(q.Get("radius"), 10, 32)
	radius = radius * 1000
	if radius == 0 {
		radius = p.Radius
	}

	maxDepartures, _ := strconv.ParseInt(q.Get("maxDepartures"), 10, 32)
	if maxDepartures == 0 {
		maxDepartures = p.MaxDepartures
	}

	maxStops, _ := strconv.ParseInt(q.Get("maxStops"), 10, 32)
	if maxStops == 0 {
		maxStops = p.MaxStops
	}

	address := q.Get("adress")
	if address != "" {
		lat, lon = geocode.GetCoordinates(address, p.Geo)
	}
	if lat == 0 || lon == 0 {
		http.Error(w, "Missing or incorrect parameters lat,lon or address", http.StatusUnprocessableEntity)
		return
	}
	res, err := GetResult(p.Repo, lat, lon, int(radius), int(maxDepartures), int(maxStops))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if len(res) > 0 {
		err = json.NewEncoder(w).Encode(res)
	} else {
		err = json.NewEncoder(w).Encode(make([]struct{}, 0))
	}
	http.Error(w, fmt.Sprintf("Error encoding JSON %s", err), http.StatusInternalServerError)
	return
}
