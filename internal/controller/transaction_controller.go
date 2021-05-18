package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"net/http"
	"time"
)

func NewTransactionController(logger logger.Logger, transactionSvc service.TransactionService, mq mq.MQ) *TransactionController {
	return &TransactionController{
		logger:         logger,
		transactionSvc: transactionSvc,
		mq:             mq,
	}
}

type TransactionController struct {
	_baseController
	mq             mq.MQ
	transactionSvc service.TransactionService
	logger         logger.Logger
}

func (c *TransactionController) CreateTransaction(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _createTransaction, func(requestID string, data []byte) (bool, error) {
		req := model.Transaction{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		transaction, err := c.transactionSvc.CreateTransaction(ctx, req)
		if err != nil {
			message.ResponseCode = http.StatusInternalServerError
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, transaction)
		if err != nil {
			return true, err
		}
		err = c.mq.Publish(message.GatewayTopic, requestID, data)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(requestID string, e error) {
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("CreateTransaction error")
	})
}

func (c *TransactionController) ListTransaction(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _listTransaction, func(requestID string, data []byte) (bool, error) {
		req := model.ListTransactionRequest{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		createGte := time.Time{}
		if req.CreatedAtGte > 0 {
			createGte = time.Unix(req.CreatedAtGte, 0)
		}
		createLte := time.Time{}
		if req.CreatedAtLte > 0 {
			createLte = time.Unix(req.CreatedAtLte, 0)
		}

		transactions, total, err := c.transactionSvc.ListTransaction(ctx, model.Transaction{
			ID:         req.ID,
			UserID:     message.User.ID,
			Type:       req.Type,
			From:       req.From,
			To:         req.To,
			CurrencyID: req.CurrencyID,
			Status:     req.Status,
			Remark:     req.Remark,
		},
			model.BaseFilter{
				CreateAtGte: createGte,
				CreateAtLte: createLte,
			}, req.Pagination, req.Sort)
		if err != nil {
			message.ResponseCode = http.StatusInternalServerError
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, model.ListTransactionResponse{
			Total:        total,
			Transactions: transactions,
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
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("ListTransaction error")
	})
}
