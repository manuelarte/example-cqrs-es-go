package api

import "github.com/google/uuid"

type GetUserByIDQuery struct {
	ID uuid.UUID
}

type GetAllUsersQuery struct {
	// TODO add pagination
}
