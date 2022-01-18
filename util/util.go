package util

import (
	"alta-wedding/config"

	"github.com/kelvins/geocoder"
)

// function untuk generate latitude, longitude menggunakan api geocode (google)
func GetGeocodeLocations(fulladdress string) (float64, float64, error) {
	env := config.GetConfig()
	apikey := env["APIKEYS"]
	geocoder.ApiKey = apikey
	var lng, lat float64
	var address geocoder.Address
	address.City = fulladdress
	location, err := geocoder.Geocoding(address)
	if err != nil {
		return 0, 0, err
	} else {
		lng = location.Latitude
		lat = location.Longitude
	}
	return lng, lat, nil
}
