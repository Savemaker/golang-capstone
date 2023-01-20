/*
Haversine formula implementation code was taken from https://github.com/umahmood/haversine
*/
package management

import "math"

type HaversineUserFinder struct {
}

const (
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns the distance in kilometers.
func (h *HaversineUserFinder) Distance(x, y *Location) float64 {
	lat1 := degreesToRadians(x.Latitude)
	lon1 := degreesToRadians(x.Longitude)
	lat2 := degreesToRadians(y.Latitude)
	lon2 := degreesToRadians(y.Longitude)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km := c * earthRaidusKm

	return km
}

// IsUserWithinRangeOfLocation evaluates if user is within radius of given location
// by calculating the distance between user location and given location using haversine formula
// and checking if distance is within radius.
func (h *HaversineUserFinder) IsUserWithinRangeOfLocation(user *User, location *Location, radius float64) bool {
	userLocation := user.UserLocation
	d := h.Distance(location, userLocation)
	return d <= radius
}
