package controller

import "bank-system-go/internal/model"

type Message struct {
	RequestID     string     `json:"request_id"`
	User          model.User `json:"user"`
	GatewayTopic  string     `json:"gateway_topic"`
	ResponseCode  int        `json:"response_code"`
	ResponseError string     `json:"response_error"`
	Payload       []byte     `json:"payload"`
}
