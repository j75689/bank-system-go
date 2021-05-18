package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletType uint8

const (
	GeneralWallet WalletType = iota + 1
)

type Wallet struct {
	ID            uint64          `json:"id" gorm:"primarykey"`
	UserID        uint64          `json:"user_id" gorm:"index"`
	Type          WalletType      `json:"type"`
	AccountNumber string          `json:"account_number" gorm:"type:varchar(255);uniqueIndex"`
	CurrencyID    uint64          `json:"currency_id"`
	Currency      Currency        `json:"currency" gorm:"->"` // read only
	Balance       decimal.Decimal `json:"balance" gorm:"not null;default:'0.0'"`
	MinWithdrawal decimal.Decimal `json:"min_withdrawal" gorm:"not null;default:'100.0'"`
	MaxWithdrawal decimal.Decimal `json:"max_withdrawal" gorm:"not null;default:'10000.0'"`
	MinDeposit    decimal.Decimal `json:"min_deposit" gorm:"not null;default:'10.0'"`
	MaxDeposit    decimal.Decimal `json:"max_deposit" gorm:"not null;default:'10000.0'"`
	MaxTransfer   decimal.Decimal `json:"max_transfer" gorm:"not null;default:'1.0'"`
	MinTransfer   decimal.Decimal `json:"min_transfer" gorm:"not null;default:'1000.0'"`
	CreatedAt     time.Time       `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt     *time.Time      `json:"-" gorm:"index"`
}

func (wallet Wallet) Preload(db *gorm.DB) *gorm.DB {
	return db.Preload("Currency")
}

type FeeType int

const (
	FIXED FeeType = iota + 1
	RATIO
)

type Currency struct {
	ID        uint64          `json:"id" gorm:"primarykey"`
	Code      string          `json:"code" gorm:"type:varchar(15)"`
	Name      string          `json:"name" gorm:"type:varchar(255)"`
	FeeType   FeeType         `json:"fee_type"`
	FeeValue  decimal.Decimal `json:"fee_value" gorm:"not null;default:'0.0'"`
	CreatedAt time.Time       `json:"created_at" gorm:"not null;default:now()"`
}

type WalletHistory struct {
	ID              uint64          `json:"id" gorm:"primarykey"`
	RequestID       string          `json:"request_id" gorm:"type:varchar(255);uniqueIndex:idx_request_account_type"`
	UserID          uint64          `json:"user_id" gorm:"index"`
	TransactionType TransactionType `json:"transaction_type" gorm:"uniqueIndex:idx_request_account_type"`
	AccountNumber   string          `json:"account_number" gorm:"type:varchar(255);uniqueIndex:idx_request_account_type"`
	Amount          decimal.Decimal `json:"amount" gorm:"not null;default:'0.0'"`
	Fee             decimal.Decimal `json:"fee" gorm:"not null;default:'0.0'"`
	CreatedAt       time.Time       `json:"created_at" gorm:"not null;default:now()"`
}
