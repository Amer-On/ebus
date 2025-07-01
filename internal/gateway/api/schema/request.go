package schema

type SubscribeRequest struct {
	Topic           string `json:"topic"`
	Event           string `json:"event"`
	CallbackAddress string `json:"callback_address"`
}
