package config

type Config struct {
	AuthToken
	UserConfig
}

type AuthToken struct {
	Token string `json:"token"`
}

type UserConfig struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangePassword struct {
	CurrentPassword string `json:"password" survey:"curr_password"`
	NewPassword string `json:"newPassword" survey:"new_password"`
}
