package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType uint8

const (
	Deposit TransactionType = iota + 1
	Withdrawal
	InternalTransfer
	ExternalTransfer
)

type TransactionStatus uint8

const (
	StatusOK TransactionStatus = iota + 1
	StatusFailed
)

type Transaction struct {
	ID         uint64            `json:"id" gorm:"primarykey"`
	UserID     uint64            `json:"user_id" gorm:"index"`
	Type       TransactionType   `json:"transaction_type"`
	From       string            `json:"from"`
	To         string            `json:"to"`
	CurrencyID uint64            `json:"currency_id"`
	Status     TransactionStatus `json:"status"`
	Amount     decimal.Decimal   `json:"amount" gorm:"not null;default:'0.0'"`
	Balance    decimal.Decimal   `json:"balance" gorm:"not null;default:'0.0'"`
	Remark     string            `json:"remark" gorm:"type:varchar(50)"`
	CreatedAt  time.Time         `json:"created_at" gorm:"not null;default:now()"`
}
