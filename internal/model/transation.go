package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransationType uint8

const (
	Deposit TransationType = iota + 1
	Withdrawal
	InternalTransfer
	ExternalTransfer
)

type Transation struct {
	ID         uint64          `json:"id" gorm:"primarykey"`
	Type       TransationType  `json:"transation_type"`
	From       string          `json:"from"`
	To         string          `json:"to"`
	CurrencyID uint64          `json:"currency_id"`
	Amount     decimal.Decimal `json:"amount" gorm:"D(20,16);not null;default:0.0"`
	Balance    decimal.Decimal `json:"balance" gorm:"decimal(20,16);not null;default:0.0"`
	Remark     string          `json:"remark" gorm:"nvarchar(50)"`
	CreatedAt  time.Time       `json:"created_at" gorm:"not null;default:now()"`
}
