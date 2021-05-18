package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"net/http"
)

func NewTransationController(logger logger.Logger, transationSvc service.TransationService, mq mq.MQ) *TransationController {
	return &TransationController{
		logger:        logger,
		transationSvc: transationSvc,
		mq:            mq,
	}
}

type TransationController struct {
	_baseController
	mq            mq.MQ
	transationSvc service.TransationService
	logger        logger.Logger
}

func (c *TransationController) CreateTransation(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _createTransation, func(requestID string, data []byte) (bool, error) {
		req := model.Transation{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		transation, err := c.transationSvc.CreateTransation(ctx, req)
		if err != nil {
			message.ResponseCode = http.StatusInternalServerError
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, transation)
		if err != nil {
			return true, err
		}
		err = c.mq.Publish(message.GatewayTopic, requestID, data)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(requestID string, e error) {
		c.logger.Error().Str("request_id", requestID).Err(e).Msg("CreateTransation error")
	})
}
