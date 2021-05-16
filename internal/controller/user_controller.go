package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"net/http"
)

func NewUserController(logger logger.Logger, userSvc service.UserService, mq mq.MQ) *UserController {
	return &UserController{
		logger:  logger,
		userSvc: userSvc,
		mq:      mq,
	}
}

type UserController struct {
	_baseController
	mq      mq.MQ
	userSvc service.UserService
	logger  logger.Logger
}

func (c *UserController) CreateUser(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _createUser, func(requestID string, data []byte) (bool, error) {
		user := model.User{}
		message, err := c.Bind(data, &user)
		if err != nil {
			return true, err
		}

		user, err = c.userSvc.Register(ctx, user)
		if err != nil {
			message.ResponseCode = http.StatusConflict
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, user)
		if err != nil {
			return true, err
		}
		return true, c.mq.Publish(message.GatewayTopic, requestID, data)
	})
}

func (c *UserController) UserLogin(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _userLogin, func(requestID string, data []byte) (bool, error) {
		req := model.UserLoginRequest{}
		message, err := c.Bind(data, &req)
		if err != nil {
			return true, err
		}

		token, err := c.userSvc.Login(ctx, req.Account, req.Password, req.IP)
		if err != nil {
			message.ResponseCode = http.StatusUnauthorized
			message.ResponseError = err.Error()
		}
		data, err = c.MarshalMessage(message, model.UserLoginResponse{
			Type:  "Bearer",
			Token: token,
		})
		if err != nil {
			return true, err
		}
		return true, c.mq.Publish(message.GatewayTopic, requestID, data)
	})
}
