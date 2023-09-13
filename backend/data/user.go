package data

type UserId string

type User struct {
	UserId   UserId `json:"client_id"`
	Username string `json:"username"`
	Device   string `json:"device"`
}

func CreateUser(Username string, Device string) *User {
	return &User{
		UserId:   "123",
		Username: Username,
		Device:   Device,
	}
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetDevice() string {
	return u.Device
}

func (u *User) UpdateUsername(Username string) {
	u.Username = Username
}
