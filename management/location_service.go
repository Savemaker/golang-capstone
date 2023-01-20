package management

type LocationService struct {
	Repository LocationRepository
	UserFinder UserFinderService
}

type LocationRepository interface {
	FindAll() []User
	UpdateUserLocation(user *User)
}

type UserFinderService interface {
	IsUserWithinRangeOfLocation(user *User, location *Location, radius float64) bool
	Distance(x, y *Location) float64
}

type User struct {
	Name         string    `json:"userName"`
	UserLocation *Location `json:"userLocation"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *LocationService) SearchUsersNearby(location *Location, radius float64) []User {
	foundUsers := make([]User, 0)
	for _, user := range l.Repository.FindAll() {
		if l.UserFinder.IsUserWithinRangeOfLocation(&user, location, radius) {
			foundUsers = append(foundUsers, user)
		}
	}
	return foundUsers
}

func (l *LocationService) UpdateUserLocation(user *User) {
	l.Repository.UpdateUserLocation(user)
}
