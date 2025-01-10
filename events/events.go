package events

type Event[T any] struct {
	Subject string `json:"subject"`
	Data    T      `json:"data"`
}
