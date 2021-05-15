package http

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Account  string `json:"account"`
	Password string `json:"password"`
}
