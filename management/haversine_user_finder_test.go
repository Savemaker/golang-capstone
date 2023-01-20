package management

import (
	"math"
	"testing"
)

func TestHaversineUserFinder(t *testing.T) {
	userFinder := HaversineUserFinder{}
	testCases := []struct {
		desription           string
		locationA, locationB Location
		expectedDistance     float64
	}{
		{
			desription:       "test distance from Moscow to Belgrade",
			locationA:        Location{Latitude: 55.750446, Longitude: 37.617494},
			locationB:        Location{Latitude: 44.817813, Longitude: 20.456897},
			expectedDistance: 1713,
		},
		{
			desription:       "test distance from Moscow to Moscow",
			locationA:        Location{Latitude: 55.750446, Longitude: 37.617494},
			locationB:        Location{Latitude: 55.750446, Longitude: 37.617494},
			expectedDistance: 0,
		},
		{
			desription:       "test distance from London to New York",
			locationA:        Location{Latitude: 51.507322, Longitude: -0.127647},
			locationB:        Location{Latitude: 40.712728, Longitude: -74.006015},
			expectedDistance: 5570,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desription, func(t *testing.T) {
			result := userFinder.Distance(&testCase.locationA, &testCase.locationB)
			roundedRes := math.Round(result)
			if roundedRes != testCase.expectedDistance {
				t.Errorf("expected %v but got %v", testCase.expectedDistance, roundedRes)
			}
		})
	}
}
