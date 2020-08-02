package geocode

import (
	geo "github.com/martinlindhe/google-geolocate"
	"log"
)

func GetCoordinates(address string, geo *geo.GoogleGeo) (lat float64, lon float64) {
	res, err := geo.Geocode(address)
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	return res.Lat, res.Lng
}
