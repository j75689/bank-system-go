package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type WalletType uint8

const (
	GeneralWallet WalletType = iota + 1
)

type Wallet struct {
	ID            uint64          `json:"id" gorm:"primarykey"`
	UserID        uint64          `json:"user_id" gorm:"index"`
	Type          WalletType      `json:"type"`
	AccountNumber string          `json:"account_number" gorm:"varchar(255);unique_index"`
	CurrencyID    uint64          `json:"currency_id"`
	Currency      Currency        `json:"currency"`
	Balance       decimal.Decimal `json:"balance" gorm:"decimal(20,16);not null;default:0.0"`
	CreatedAt     time.Time       `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt     *time.Time      `json:"deleted_at" gorm:"index"`
}

type Currency struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	Code      string    `json:"code" gorm:"varchar(15)"`
	Name      string    `json:"name" gorm:"nvarchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}
