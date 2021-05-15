package controller

type Message struct {
	RequestID    string `json:"reuest_id"`
	GatewayTopic string `json:"gateway_topic"`
	ResponseCode int
	Payload      []byte `json:"payload"`
}
