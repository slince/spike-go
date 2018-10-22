package auth

type Authentication interface {
	// 验证客户端信息
	Auth(credentials map[string]interface{}) error
}
