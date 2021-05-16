package model

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	User
	Password interface{} `json:"password,omitempty"` // ignore password
}

type UserLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

type UserLoginResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
