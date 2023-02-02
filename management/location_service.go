package management

type LocationService struct {
	Repository LocationRepository
	UserFinder UserFinderService
}

type LocationRepository interface {
	FindAll() ([]User, error)
	UpdateUserLocation(user *User) error
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

func (l *LocationService) SearchUsersNearby(location *Location, radius float64) ([]User, error) {
	foundUsers := make([]User, 0)
	users, err := l.Repository.FindAll()
	if err != nil {
		return foundUsers, err
	}
	for _, user := range users {
		if l.UserFinder.IsUserWithinRangeOfLocation(&user, location, radius) {
			foundUsers = append(foundUsers, user)
		}
	}
	return foundUsers, nil
}

func (l *LocationService) UpdateUserLocation(user *User) error {
	return l.Repository.UpdateUserLocation(user)
}
