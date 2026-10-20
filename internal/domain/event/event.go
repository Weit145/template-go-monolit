package event

import "github.com/google/uuid"

type Event struct {
	ID   uuid.UUID
	Type string
	Data any
}

type Command struct {
	ID   uuid.UUID
	Type string
}
