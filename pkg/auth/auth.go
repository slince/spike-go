package auth

import (
	"github.com/slince/spike/pkg/msg"
)

type Auth interface {
	Check(login *msg.Login) User
}

type SimpleAuth struct {
	Users []*GenericUser
	a     string
}

func (au *SimpleAuth) Check(login *msg.Login) User {
	for _, u := range au.Users {
		if u.Password == login.Password && u.Username == login.Username {
			return u
		}
	}
	return nil
}

func NewSimpleAuth(users []*GenericUser) *SimpleAuth {
	return &SimpleAuth{Users: users}
}
