package config

type Config struct {
	AuthToken
	UserConfig
}

type AppConfig struct {
	ProjectName string `json:"project_name"`
	Runtime string `json:"runtime"`
	Port string `json:"port"`
}

type ProjectMetaData struct {
	ProjectName string `json:"project_name"`
	CurrentVersion string `json:"current_version"`
	Versions []string `json:"versions"`
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
