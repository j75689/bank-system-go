package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"net/http"
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

		wallet, err := c.walletSvc.CreateWallet(ctx, req)
		if err != nil {
			message.ResponseCode = http.StatusUnauthorized
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, wallet)
		if err != nil {
			return true, err
		}
		return true, c.mq.Publish(message.GatewayTopic, requestID, data)
	})
}
