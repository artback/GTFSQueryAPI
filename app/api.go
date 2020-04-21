package app

import (
	"github.com/artback/gtfsQueryGoApi/config"
	"github.com/artback/gtfsQueryGoApi/direction"
	"github.com/artback/gtfsQueryGoApi/query"
	"github.com/gorilla/mux"
	geo "github.com/martinlindhe/google-geolocate"
	"log"
	"net/http"
)

var (
	conf      *config.Configuration
	repo      *query.Repository
	geoclient *geo.GoogleGeo
)

func init() {
	conf = new(config.Configuration)
	repo = new(query.Repository)
	geoclient = geo.NewGoogleGeo("api-key")
	err := config.Init(conf)
	if err != nil {
		panic(err)
	}
}

func Run() {
	err := repo.Connect(conf.Database)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/departures/place", placeHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func placeHandler(w http.ResponseWriter, r *http.Request) {
	direction.PlaceHandler(repo, w, r, conf.Default, geo)
}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
