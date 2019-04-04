package questions

type UserInfo struct {
	Name     string `json:"name" survey:"name"`
	Email    string `json:"email" survey:"email"`
	Password string `json:"password" survey:"password"`
	Random_String string `json:"random_string" survey:"random"`
}

type AuthCredential struct {
	Email    string `json:"email" survey:"email"`
	Password string `json:"password" survey:"password"`
}
