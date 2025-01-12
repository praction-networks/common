package events

type Event[T any] struct {
	Subject Subject `json:"subject"`
	Data    T       `json:"data"`
}
