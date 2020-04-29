package geocode

import (
	geo "github.com/martinlindhe/google-geolocate"
	"log"
)

func GetCordinatesForAddress(address string, geo *geo.GoogleGeo) (lat float64, lon float64) {
	res, err := geo.Geocode(address)
	if err != nil {
		log.Println(err)
	}
	return res.Lat, res.Lng
}
