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

type TransationStatus uint8

const (
	StatusOK TransationStatus = iota + 1
	StatusFailed
)

type Transation struct {
	ID         uint64           `json:"id" gorm:"primarykey"`
	UserID     uint64           `json:"user_id" gorm:"index"`
	Type       TransationType   `json:"transation_type"`
	From       string           `json:"from"`
	To         string           `json:"to"`
	CurrencyID uint64           `json:"currency_id"`
	Status     TransationStatus `json:"status"`
	Amount     decimal.Decimal  `json:"amount" gorm:"not null;default:'0.0'"`
	Balance    decimal.Decimal  `json:"balance" gorm:"not null;default:'0.0'"`
	Remark     string           `json:"remark" gorm:"type:varchar(50)"`
	CreatedAt  time.Time        `json:"created_at" gorm:"not null;default:now()"`
}
