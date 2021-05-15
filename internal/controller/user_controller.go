package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/service"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"encoding/json"
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
	callbackChans _callbackMap
	mq            mq.MQ
	userSvc       service.UserService
	logger        logger.Logger
}

func (c *UserController) CreateUser(ctx context.Context) error {
	return c.mq.Subscribe(ctx, _createUser, func(requestID string, data []byte) (bool, error) {
		c.logger.Info().Str("request_id", requestID).Bytes("message", data).Send()
		message := Message{
			ResponseCode: http.StatusOK,
		}
		err := json.Unmarshal(data, &message)
		if err != nil {
			return true, err
		}

		user := model.User{}
		err = json.Unmarshal(message.Payload, &user)
		if err != nil {
			return true, err
		}

		user, err = c.userSvc.Register(ctx, user)
		payload := []byte{}
		if err != nil {
			message.ResponseCode = http.StatusConflict
			payload = []byte(err.Error())
		} else {
			payload, err = json.Marshal(user)
			if err != nil {
				return true, err
			}
		}
		message.Payload = payload
		data, err = json.Marshal(message)
		if err != nil {
			return true, err
		}
		return true, c.mq.Publish(message.GatewayTopic, requestID, data)
	})
}
