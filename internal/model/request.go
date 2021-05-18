package model

import (
	"github.com/shopspring/decimal"
)

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

type ListWalletRequest struct {
	Type          WalletType `json:"type"`
	AccountNumber string     `json:"account_number"`
	CurrencyID    uint64     `json:"currency_id"`
	Sort          Sorting    `json:"sort"`
	Pagination    Pagination `json:"pagination"`
}

type ListWalletResponse struct {
	Total   int64    `json:"total"`
	Wallets []Wallet `json:"wallets"`
}

type UpdateWalletBalanceRequest struct {
	Type          TransactionType `json:"type"`
	AccountNumber string          `json:"account_number"`
	Amount        decimal.Decimal `json:"amount"`
}

type UpdateWalletBalanceResponse struct {
	Error      *string      `json:"error,omitempty"`
	Wallet     *Wallet      `json:"wallet,omitempty"`
	Transation *Transaction `json:"transation,omitempty"`
}

type ListTransactionRequest struct {
	ID           uint64            `json:"id" gorm:"primarykey"`
	Type         TransactionType   `json:"transaction_type"`
	From         string            `json:"from"`
	To           string            `json:"to"`
	CurrencyID   uint64            `json:"currency_id"`
	Status       TransactionStatus `json:"status"`
	Remark       string            `json:"remark" gorm:"type:varchar(50)"`
	CreatedAtGte int64             `json:"created_at_gte" gorm:"not null;default:now()"`
	CreatedAtLte int64             `json:"created_at_lte" gorm:"not null;default:now()"`
	Sort         Sorting           `json:"sort"`
	Pagination   Pagination        `json:"pagination"`
}

type ListTransactionResponse struct {
	Total        int64         `json:"total"`
	Transactions []Transaction `json:"transactions"`
}

type TransferRequest struct {
	Type   TransactionType `json:"type"`
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

type TransferResponse struct {
	Error      *string      `json:"error,omitempty"`
	Wallet     *Wallet      `json:"wallet,omitempty"`
	Transation *Transaction `json:"transation,omitempty"`
}
