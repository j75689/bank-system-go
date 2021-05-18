package error

import "errors"

var (
	ErrWalletAccountNotFound          = errors.New("account not found")
	ErrWalletBalanceInsufficient      = errors.New("account balance is insufficient")
	ErrLessThanMinDepositAmount       = errors.New("amount is less than minimin of deposit")
	ErrGreaterThanMaxDepositAmount    = errors.New("amount is greater than maximin of deposit")
	ErrLessThanMinWithdrawalAmount    = errors.New("amount is less than minimin of withdrawal")
	ErrGreaterThanMaxWithdrawalAmount = errors.New("amount is greater than maximin of withdrawal")
	ErrLessThanMinTransferAmount      = errors.New("amount is less than minimin of transfer")
	ErrGreaterThanMaxTransferAmount   = errors.New("amount is greater than maximin of transfer")
)
