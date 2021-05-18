package controller

import (
	pkgErr "bank-system-go/internal/errors"
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"
)

func NewWalletController(logger logger.Logger, walletSvc service.WalletService, mq mq.MQ) *WalletController {
	return &WalletController{
		logger:    logger,
		walletSvc: walletSvc,
		mq:        mq,
	}
}

type WalletController struct {
	_baseController
	mq        mq.MQ
	walletSvc service.WalletService
	logger    logger.Logger
}

func (c *WalletController) CreateWallet(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _createWallet, func(requestID string, data []byte) (bool, error) {
		req := model.Wallet{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		// default
		{
			req.MaxDeposit = decimal.NewFromInt(10000)
			req.MaxWithdrawal = decimal.NewFromInt(10000)
			req.MaxTransfer = decimal.NewFromInt(1000)
			req.MinDeposit = decimal.NewFromInt(10)
			req.MinWithdrawal = decimal.NewFromInt(100)
			req.MinTransfer = decimal.NewFromFloat(0.10)
		}

		wallet, err := c.walletSvc.CreateWallet(ctx, req)
		if err != nil {
			message.ResponseCode = http.StatusInternalServerError
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, wallet)
		if err != nil {
			return true, err
		}
		err = c.mq.Publish(message.GatewayTopic, requestID, data)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(requestID string, e error) {
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("CreateWallet error")
	})
}

func (c *WalletController) ListWallet(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _listWallet, func(requestID string, data []byte) (bool, error) {
		req := model.ListWalletRequest{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		wallets, total, err := c.walletSvc.ListWallet(ctx, model.Wallet{
			UserID:        message.User.ID,
			CurrencyID:    req.CurrencyID,
			Type:          req.Type,
			AccountNumber: req.AccountNumber,
		}, req.Pagination, req.Sort)
		if err != nil {
			message.ResponseCode = http.StatusInternalServerError
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, &model.ListWalletResponse{
			Total:   total,
			Wallets: wallets,
		})
		if err != nil {
			return true, err
		}
		err = c.mq.Publish(message.GatewayTopic, requestID, data)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(requestID string, e error) {
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("ListWallet error")
	})
}

func (c *WalletController) checkMinMax(transactionType model.TransactionType, wallet model.Wallet, amount decimal.Decimal) error {
	switch transactionType {
	case model.Deposit:
		if amount.LessThan(wallet.MinDeposit) {
			return pkgErr.ErrLessThanMinDepositAmount
		}
		if amount.GreaterThan(wallet.MaxDeposit) {
			return pkgErr.ErrGreaterThanMaxDepositAmount
		}
	case model.Withdrawal:
		if amount.LessThan(wallet.MinWithdrawal) {
			return pkgErr.ErrLessThanMinWithdrawalAmount
		}
		if amount.GreaterThan(wallet.MaxWithdrawal) {
			return pkgErr.ErrGreaterThanMaxWithdrawalAmount
		}
	}
	return nil
}

func (c *WalletController) UpdateWalletBalance(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _updateWalletBalance, func(requestID string, data []byte) (ack bool, err error) {
		req := model.UpdateWalletBalanceRequest{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}
		ack = true
		resp := &model.UpdateWalletBalanceResponse{}

		defer func() {
			data, err = c.MarshalMessage(message, resp)
			if err != nil {
				ack = true
			}

			err = c.mq.Publish(message.GatewayTopic, requestID, data)
			if err != nil {
				ack = false
			}
		}()

		filter := model.Wallet{
			UserID:        message.User.ID,
			AccountNumber: req.AccountNumber,
		}
		wallet, err := c.walletSvc.GetWallet(ctx, filter)
		if err != nil {
			if errors.Is(err, pkgErr.ErrWalletAccountNotFound) {
				msg := err.Error()
				resp.Error = &msg
			} else {
				message.ResponseCode = http.StatusInternalServerError
				message.ResponseError = err.Error()
			}
			return
		}
		resp.Wallet = &wallet

		err = c.checkMinMax(req.Type, wallet, req.Amount.Abs())
		if err != nil {
			msg := err.Error()
			resp.Error = &msg
			return
		}

		status := model.StatusOK
		err = c.walletSvc.UpdateBalance(ctx, filter, requestID, req.Type, req.Amount)
		if err != nil {
			if errors.Is(err, pkgErr.ErrWalletBalanceInsufficient) {
				msg := err.Error()
				resp.Error = &msg
			} else {
				message.ResponseCode = http.StatusInternalServerError
				message.ResponseError = err.Error()
			}
			status = model.StatusFailed
		}
		wallet, err = c.walletSvc.GetWallet(ctx, filter)
		if err != nil {
			if errors.Is(err, pkgErr.ErrWalletAccountNotFound) {
				msg := err.Error()
				resp.Error = &msg
			} else {
				message.ResponseCode = http.StatusInternalServerError
				message.ResponseError = err.Error()
			}
			return
		}
		resp.Wallet = &wallet

		// TODO: outbox pattern
		{
			data, err = c.MarshalMessage(message, &model.Transaction{
				UserID:     message.User.ID,
				Type:       req.Type,
				Status:     status,
				From:       wallet.AccountNumber,
				To:         wallet.AccountNumber,
				CurrencyID: wallet.CurrencyID,
				Amount:     req.Amount.Abs(),
				Balance:    wallet.Balance,
			})
			if err != nil {
				return true, err
			}
			err = c.mq.Publish(_createTransaction, requestID, data)
			if err != nil {
				return false, err
			}
		}

		return
	}, func(requestID string, e error) {
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("UpdateWalletBalance error")
	})
}
