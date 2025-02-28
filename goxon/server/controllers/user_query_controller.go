package controllers

import (
	"codeberg.org/manuelarte/goxon/pkg/query"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"main/api"
	"time"
)

type UserQueryController struct {
	QueryGateway query.QueryGateway
}

func (u *UserQueryController) HandleGetUserByIDQuery(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fiber.StatusBadRequest, "ID can't be empty")
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ID isn't valid",
		})
	}

	q := api.GetUserByIDQuery{ID: userId}
	ctx, cancel := context.WithDeadline(c.Context(), time.Now().Add(300*time.Millisecond))
	defer cancel()
	ch, err := u.QueryGateway.Query(ctx, q)
	if err != nil {
		return err
	}
	defer close(ch)
	select {
	case result := <-ch:
		if result.Err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Query error: " + result.Err.Error(),
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(result.Data)
		}
	case <-ctx.Done():
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Query error timeout",
		})
	}

}
