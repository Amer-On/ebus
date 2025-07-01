package domain

type Subscription struct {
	Subscriber
	Topic string `json:"topic"`
	Event string `json:"event"`
}

type Subscriber struct {
	CallbackAddress string `json:"callback_address"`
	Name            string `json:"name"`
}
