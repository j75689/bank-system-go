package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type _callbackMap struct {
	lock  sync.RWMutex
	chans map[string]chan []byte
}

func (m *_callbackMap) Get(key string) (chan []byte, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	c, ok := m.chans[key]
	return c, ok
}

func (m *_callbackMap) Set(key string, value chan []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.chans[key] = value
}

func (m *_callbackMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.chans, key)
}

func NewGatewayController(logger logger.Logger, mq mq.MQ) *GatewayController {
	return &GatewayController{
		callbackChans: _callbackMap{
			lock:  sync.RWMutex{},
			chans: make(map[string]chan []byte),
		},
		logger: logger,
		mq:     mq,
	}
}

type GatewayController struct {
	callbackChans _callbackMap
	mq            mq.MQ
	logger        logger.Logger
}

func (c *GatewayController) RegisterUser(ctx context.Context, requestID, gatewayTopic string, user model.User) (int, model.User, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return http.StatusInternalServerError, model.User{}, err
	}
	message := Message{
		RequestID:    requestID,
		GatewayTopic: gatewayTopic,
		Payload:      payload,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return http.StatusInternalServerError, model.User{}, err
	}
	err = c.mq.Publish(_createUser, requestID, data)
	if err != nil {
		return http.StatusInternalServerError, model.User{}, err
	}

	callback := make(chan []byte)
	c.callbackChans.Set(requestID, callback)
	data = <-callback
	close(callback)
	c.callbackChans.Delete(requestID)

	err = json.Unmarshal(data, &message)
	if err != nil {
		return http.StatusInternalServerError, model.User{}, err
	}
	err = json.Unmarshal(message.Payload, &user)
	if err != nil {
		return message.ResponseCode, model.User{}, errors.WithMessage(errors.New(string(message.Payload)), "create uer failed")
	}
	return message.ResponseCode, user, nil
}

func (c *GatewayController) GatewayCallback(ctx context.Context, gatewayTopic string) error {
	return c.mq.Subscribe(ctx, gatewayTopic, func(requestID string, data []byte) (bool, error) {
		callback, ok := c.callbackChans.Get(requestID)
		if ok {
			callback <- data
		}
		return true, nil
	})
}
