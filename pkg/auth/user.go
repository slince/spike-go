package auth

type User interface {
	getUsername() string
}

type GenericUser struct {
	Username string
	Password string
}

func (u *GenericUser) getUsername() string {
	return u.Username
}
