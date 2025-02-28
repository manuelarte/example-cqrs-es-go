package api

import (
	"github.com/google/uuid"
	"time"
)

type UserCreatedEvent struct {
	ID        uuid.UUID `goxon:"targetAggregateID"`
	Name      string
	Surname   string
	BirthDate time.Time
}
