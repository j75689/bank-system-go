package controller

import (
	"bank-system-go/internal/model"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
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

func NewGatewayController(logger logger.Logger, mq mq.MQ) (*GatewayController, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &GatewayController{
		callbackTopic: hostname,
		callbackChans: _callbackMap{
			lock:  sync.RWMutex{},
			chans: make(map[string]chan []byte),
		},
		logger: logger,
		mq:     mq,
	}, nil
}

type GatewayController struct {
	_baseController
	callbackTopic string
	callbackChans _callbackMap
	mq            mq.MQ
	logger        logger.Logger
}

func (c *GatewayController) Wait(requestID string) []byte {
	callback := make(chan []byte)
	c.callbackChans.Set(requestID, callback)
	data := <-callback
	close(callback)
	c.callbackChans.Delete(requestID)
	return data
}

func (c *GatewayController) GatewayCallback(ctx context.Context) error {
	return c.mq.Subscribe(ctx, c.callbackTopic, func(requestID string, data []byte) (bool, error) {
		callback, ok := c.callbackChans.Get(requestID)
		if ok {
			callback <- data
		}
		return true, nil
	})
}

func (c *GatewayController) PushMessage(requestID string, code int, topic string, user model.User, model interface{}) error {
	payload, err := json.Marshal(model)
	if err != nil {
		return err
	}
	message := Message{
		RequestID:    requestID,
		User:         user,
		GatewayTopic: c.callbackTopic,
		ResponseCode: code,
		Payload:      payload,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.mq.Publish(topic, requestID, data)
}

func (c *GatewayController) RegisterUser(ctx context.Context, requestID string, req model.RegisterUserRequest) (int, model.RegisterUserResponse, error) {
	user := model.User{
		Name:     req.Name,
		Account:  req.Account,
		Password: []byte(req.Password),
	}
	err := c.PushMessage(requestID, http.StatusOK, _createUser, model.User{}, user)
	if err != nil {
		return http.StatusInternalServerError, model.RegisterUserResponse{}, err
	}
	data := c.Wait(requestID)

	message, err := c.Bind(data, &user)
	if err != nil {
		return message.ResponseCode, model.RegisterUserResponse{}, errors.WithMessage(err, "create user error")
	}
	return message.ResponseCode, model.RegisterUserResponse{
		User: user,
	}, nil
}

func (c *GatewayController) Login(ctx context.Context, requestID string, req model.UserLoginRequest) (int, model.UserLoginResponse, error) {
	err := c.PushMessage(requestID, http.StatusOK, _userLogin, model.User{}, req)
	if err != nil {
		return http.StatusInternalServerError, model.UserLoginResponse{}, err
	}

	data := c.Wait(requestID)

	resp := model.UserLoginResponse{}
	message, err := c.Bind(data, &resp)
	if err != nil {
		return message.ResponseCode, model.UserLoginResponse{}, errors.WithMessage(err, "user login error")
	}
	return message.ResponseCode, resp, nil
}

func (c *GatewayController) VerifyUser(ctx context.Context, requestID string, req model.VerifyUserRequest) (int, model.VerifyUserResponse, error) {
	err := c.PushMessage(requestID, http.StatusOK, _verifyUser, model.User{}, req)
	if err != nil {
		return http.StatusInternalServerError, model.VerifyUserResponse{}, err
	}

	data := c.Wait(requestID)

	resp := model.VerifyUserResponse{}
	message, err := c.Bind(data, &resp)
	if err != nil {
		return message.ResponseCode, model.VerifyUserResponse{}, errors.WithMessage(err, "verify user error")
	}
	return message.ResponseCode, resp, nil
}

func (c *GatewayController) CreateWallet(ctx context.Context, requestID string, user model.User, req model.CreateWalletRequest) (int, model.CreateWalletResponse, error) {
	wallet := model.Wallet{
		UserID:     user.ID,
		Type:       req.Type,
		CurrencyID: req.CurrencyID,
	}
	err := c.PushMessage(requestID, http.StatusOK, _createWallet, user, wallet)
	if err != nil {
		return http.StatusInternalServerError, model.CreateWalletResponse{}, err
	}

	data := c.Wait(requestID)

	message, err := c.Bind(data, &wallet)
	if err != nil {
		return message.ResponseCode, model.CreateWalletResponse{}, errors.WithMessage(err, "create wallet error")
	}
	return message.ResponseCode, model.CreateWalletResponse{
		Wallet: wallet,
	}, nil
}

func (c *GatewayController) ListWallet(ctx context.Context, requestID string, user model.User, req model.ListWalletRequest) (int, model.ListWalletResponse, error) {
	err := c.PushMessage(requestID, http.StatusOK, _listWallet, user, req)
	if err != nil {
		return http.StatusInternalServerError, model.ListWalletResponse{}, err
	}

	data := c.Wait(requestID)

	resp := model.ListWalletResponse{}
	message, err := c.Bind(data, &resp)
	if err != nil {
		return message.ResponseCode, model.ListWalletResponse{}, errors.WithMessage(err, "list wallet error")
	}
	return message.ResponseCode, resp, nil
}

func (c *GatewayController) UpdateWalletBalance(ctx context.Context, requestID string, user model.User, req model.UpdateWalletBalanceRequest) (int, model.UpdateWalletBalanceResponse, error) {
	if !(req.Type == model.Withdrawal || req.Type == model.Deposit) {
		return http.StatusBadRequest, model.UpdateWalletBalanceResponse{}, errors.New("type is not withdrawal or deposit")
	}
	if req.Type == model.Withdrawal {
		req.Amount = req.Amount.Abs().Mul(decimal.NewFromInt(-1))
	}
	if req.Type == model.Deposit {
		req.Amount = req.Amount.Abs()
	}
	err := c.PushMessage(requestID, http.StatusOK, _updateWalletBalance, user, req)
	if err != nil {
		return http.StatusInternalServerError, model.UpdateWalletBalanceResponse{}, err
	}

	data := c.Wait(requestID)

	resp := model.UpdateWalletBalanceResponse{}
	message, err := c.Bind(data, &resp)
	if err != nil {
		return message.ResponseCode, model.UpdateWalletBalanceResponse{}, errors.WithMessage(err, "list wallet error")
	}
	return message.ResponseCode, resp, nil
}
