package types

import (
	"net/http"
)

type UserInfo struct {
	Name     string `json:"name" survey:"name"`
	Email    string `json:"email" survey:"email"`
	Password string `json:"password" survey:"password"`
}

type AuthCredential struct {
	Email    string `json:"email" survey:"email"`
	Password string `json:"password" survey:"password"`
}

type HttpClient struct {
	Client *http.Client
}
