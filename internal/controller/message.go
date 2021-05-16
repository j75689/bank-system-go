package controller

type Message struct {
	RequestID     string `json:"reuest_id"`
	GatewayTopic  string `json:"gateway_topic"`
	ResponseCode  int    `json:"response_code"`
	ResponseError string `json:"response_error"`
	Payload       []byte `json:"payload"`
}
