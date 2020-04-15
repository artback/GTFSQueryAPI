package direction

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func PlaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lat, _ := strconv.ParseFloat(vars["lat"], 32)
	lon, _ := strconv.ParseFloat(vars["lon"], 32)
	adress := vars["adress"]
	if adress != "" {
		lat, lon = getCordinatesForAdress(adress)
	}
}
func getCordinatesForAdress(adress string) (lat float64, lon float64) {

}
