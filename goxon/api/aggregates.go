package api

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `goxon:"aggregateID"`
	Name      string
	Surname   string
	BirthDate time.Time
}
