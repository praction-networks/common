package events

type Event[T any] struct {
	Subject Subjects `json:"subject"`
	Data    T        `json:"data"`
}
