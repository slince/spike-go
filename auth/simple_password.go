package auth

import "fmt"

type SimplePassword struct{
	Username string
	Password string
}
// 简单验证
func (sp *SimplePassword) Auth(credentials map[string]interface{}) error{
	username,usOk := credentials["username"]
	password,psOk := credentials["password"]

	if !usOk || !psOk {
		return fmt.Errorf("missing username or password")
	}
	if username != sp.Username || password != sp.Password {
		return fmt.Errorf("bad username or password")
	}
	return nil
}
