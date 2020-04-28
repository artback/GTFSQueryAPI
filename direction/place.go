package direction

import (
	"encoding/json"
	"github.com/artback/gtfsQueryGoApi/config"
	"github.com/artback/gtfsQueryGoApi/query"
	geo "github.com/martinlindhe/google-geolocate"
	"log"
	"net/http"
	"os"
	"strconv"
)

func PlaceHandler(repo *query.Repository, w http.ResponseWriter, r *http.Request, d config.DefaultConfiguration, geo *geo.GoogleGeo) {
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
		radius = d.Radius
	}

	maxDepartures, _ := strconv.ParseInt(q.Get("maxDepartures"), 10, 32)
	if maxDepartures == 0 {
		maxDepartures = d.MaxDepartures
	}

	maxStops, _ := strconv.ParseInt(q.Get("maxStops"), 10, 32)
	if maxStops == 0 {
		maxStops = d.MaxStops
	}

	address := q.Get("adress")
	if address != "" {
		lat, lon = getCordinatesForAddress(address, geo)
	}
	if address == "" && (lat == 0 || lon == 0) {
		http.Error(w, "Missing obligatory parameters lat,lon or adress", http.StatusUnprocessableEntity)
		return
	}
	res := GetResult(repo, lat, lon, int(radius), int(maxDepartures), int(maxStops))
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
	}
}
func getCordinatesForAddress(address string, geo *geo.GoogleGeo) (lat float64, lon float64) {
	res, _ := geo.Geocode(address)
	return res.Lat, res.Lng
}
