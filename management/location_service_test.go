package management

import "testing"

func TestLocationService(t *testing.T) {

	users := []User{
		{Name: "John", UserLocation: &Location{Latitude: 44.842465, Longitude: 20.380679}},
		{Name: "Bob", UserLocation: &Location{Latitude: 55.751197, Longitude: 37.786482}},
		{Name: "Don", UserLocation: &Location{Latitude: 55.740917, Longitude: 37.627344}},
		{Name: "Alex", UserLocation: &Location{Latitude: 55.755532, Longitude: 37.633587}},
	}

	locationService := &LocationService{Repository: &InMemoryLocationDB{Users: users}, UserFinder: &HaversineUserFinder{}}

	testCases := []struct {
		description        string
		location           *Location
		radius             float64
		expectedUsersCount int
	}{
		{
			description:        "one user should be found",
			location:           &Location{Latitude: 44.845468, Longitude: 20.410705},
			radius:             3,
			expectedUsersCount: 1,
		},
		{
			description:        "zero users should be found",
			location:           &Location{Latitude: 45.255134, Longitude: 19.845176},
			radius:             3,
			expectedUsersCount: 0,
		},
		{
			description:        "three users should be found",
			location:           &Location{Latitude: 55.7571978, Longitude: 37.659332},
			radius:             10,
			expectedUsersCount: 3,
		},
		{
			description:        "two users should be found",
			location:           &Location{Latitude: 55.7571978, Longitude: 37.659332},
			radius:             3,
			expectedUsersCount: 2,
		},
		{
			description:        "four users should be found",
			location:           &Location{Latitude: 55.7571978, Longitude: 37.659332},
			radius:             2000,
			expectedUsersCount: 4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			foundUsers, _ := locationService.SearchUsersNearby(testCase.location, testCase.radius)
			if len(foundUsers) != testCase.expectedUsersCount {
				t.Errorf("expected %v but found %v", testCase.expectedUsersCount, len(foundUsers))
			}
		})
	}
}
