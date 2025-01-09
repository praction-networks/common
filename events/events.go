package events

import eventSubjects "github.com/praction-networks/common/events/eventsubjects"

type Event[T any] struct {
	Version int                    `json:"version"`
	Subject eventSubjects.Subjects `json:"subject"`
	Data    T                      `json:"data"`
}
