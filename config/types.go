package config

type Config struct {
	AuthToken	
	UserConfig
}

type AuthToken struct {
	Token string `json:"token"`
}

type UserConfig struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Email string `json:"email"`
}