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

type VerifyUserRequest struct {
	Token string `json:"token"`
}

type VerifyUserResponse struct {
	Islegal bool `json:"is_legal"`
	User    User `json:"User"`
}

type CreateWalletRequest struct {
	Type       WalletType `json:"type"`
	CurrencyID uint64     `json:"currency_id"`
}

type CreateWalletResponse struct {
	Wallet
}
