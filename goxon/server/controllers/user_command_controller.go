package controllers

import (
	"codeberg.org/manuelarte/goxon/command"
	"codeberg.org/manuelarte/goxon/pkg/query"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"main/api"
	"time"
)

type UserCommandController struct {
	QueryGateway   query.QueryGateway
	CommandGateway command.CommandGateway
}

func (uc *UserCommandController) HandleCreate(c *fiber.Ctx) error {
	user := &api.User{}
	if err := json.Unmarshal(c.Body(), user); err != nil {
		return err
	}
	ctx, cancel := context.WithDeadline(c.Context(), time.Now().Add(300*time.Millisecond))
	defer cancel()
	qCh, err := uc.QueryGateway.Query(ctx, api.GetUserByIDQuery{ID: user.ID})
	if err != nil {
		return err
	}
	select {
	case q := <-qCh:
		if q.Data != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("user with id %s already exists", user.ID)})
		}
	case <-ctx.Done():
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New("context deadline exceeded")})
	}

	err = uc.CommandGateway.SubscribeAggregate(c.Context(), user, func() any {
		return api.UserCreatedEvent{
			ID:        user.ID,
			Name:      user.Name,
			Surname:   user.Surname,
			BirthDate: user.BirthDate,
		}
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(user)

}
