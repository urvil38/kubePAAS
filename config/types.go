package config

import (
	"github.com/urvil38/kubepaas/schema/latest"
)

type AuthConfig struct {
	AuthToken
	UserConfig
}

type AppConfig struct {
	ProjectName string `json:"project_name"`
	Runtime     string `json:"runtime"`
	Port        string `json:"port"`
	StaticDir   string `json:"static_dir"`
}

type secret string

type ProjectMetaData struct {
	ProjectName         string   `json:"project_name"`
	CurrentVersion      string   `json:"current_version"`
	Versions            []string `json:"versions"`
	GCPProject          string   `json:"gcp_project"`
	Domain              string   `json:"domain_name"`
	SourceCodeBucket    string   `json:"source_bucket"`
	CloudBuildLogBucket string   `json:"cloudbuild_bucket"`
	CloudBuildSecret    secret   `json:"cloudbuild_secret,omitempty"`
	CloudStorageSecret  secret   `json:"cloudstorage_secret,omitempty"`
}

func (secret) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

type Kubernetes struct {
	ProjectName    string                `json:"project_name"`
	CurrentVersion string                `json:"current_version"`
	Spec           latest.KubepaasConfig `json:"spec"`
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
	NewPassword     string `json:"newPassword" survey:"new_password"`
}

type KubepaasConfig struct {
	ProjectRoot  string
	KubepaasRoot string
}
