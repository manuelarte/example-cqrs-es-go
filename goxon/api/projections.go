package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type UserRead struct {
	ID   uuid.UUID
	Name string
}

type UserReadProjection struct {
	Users []*UserRead
}

func (p *UserReadProjection) HandleGetUserByIDQuery(_ context.Context, q GetUserByIDQuery) (*UserRead, error) {
	indexed := lo.KeyBy(p.Users, func(item *UserRead) uuid.UUID {
		return item.ID
	})
	if u, ok := indexed[q.ID]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user %s not found", q.ID)
}

func (p *UserReadProjection) OnUserCreatedEvent(_ context.Context, e UserCreatedEvent) {
	p.Users = append(p.Users, &UserRead{
		ID:   e.ID,
		Name: e.Name,
	})
}
