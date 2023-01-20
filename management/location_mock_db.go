package management

type InMemoryLocationDB struct {
	Users []User
}

func (i *InMemoryLocationDB) FindAll() []User {
	return i.Users
}

func (i *InMemoryLocationDB) FindUserByName(name string) *User {
	return &User{Name: "Name", UserLocation: &Location{Latitude: 0.0, Longitude: 0.1}}
}

func (i *InMemoryLocationDB) UpdateUserLocation(user *User) {
	i.Users = append(i.Users, *user)
}
