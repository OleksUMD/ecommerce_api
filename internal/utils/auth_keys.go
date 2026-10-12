package utils

type UserAuthKey string

const (
	UserRoleKey  UserAuthKey = "user_role"
	UserEmailKey UserAuthKey = "user_email"
	UserIDKey    UserAuthKey = "user_id"
	AdminRole    UserAuthKey = "admin"
)
