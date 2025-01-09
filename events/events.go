package events

type Event[T any] struct {
	Version int      `json:"version"`
	Subject Subjects `json:"subject"`
	Data    T        `json:"data"`
}
