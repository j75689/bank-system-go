package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type _baseController struct {
}

func (_baseController) Bind(data []byte, model interface{}) (Message, error) {
	message := Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		return Message{
			ResponseCode: http.StatusInternalServerError,
		}, err
	}
	if len(message.ResponseError) > 0 {
		return message, errors.New(message.ResponseError)
	}
	err = json.Unmarshal(message.Payload, model)
	if err != nil {
		return message, err
	}
	return message, nil
}

func (_baseController) MarshalMessage(message Message, model interface{}) ([]byte, error) {
	payload, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	message.Payload = payload
	return json.Marshal(message)
}
