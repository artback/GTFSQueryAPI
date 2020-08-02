package app

import (
	"github.com/allbin/gtfsQueryGoApi/direction"
	"github.com/allbin/gtfsQueryGoApi/internal/config"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/gorilla/mux"
	geo "github.com/martinlindhe/google-geolocate"
	"log"
	"net/http"
	"os"
)

func Run() {
	conf, err := config.NewConfig()
	repo, err := query.NewConnectedRepository(conf.Database)
	if err != nil {
		panic(err)
	}
	geoClient := geo.NewGoogleGeo(os.Getenv("GOOGLE_GEOCODE_API_KEY"))

	log.Print("Server up and running...")
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	places := direction.Places{Repo: repo, DefaultConfiguration: conf.Default, Geo: geoClient}
	r.HandleFunc("/departures/place", places.PlaceHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
