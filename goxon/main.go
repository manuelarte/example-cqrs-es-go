package main

import (
	"codeberg.org/manuelarte/goxon"
	"encoding/gob"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"log"
	"main/api"
	"main/server/controllers"
)

func main() {
	gob.Register(api.UserRead{})
	up := &api.UserReadProjection{}

	gf, err := goxon.NewGoxonFramework(&goxon.Config{
		Projections: []any{up},
	})
	if err != nil {
		log.Fatal(err)
	}

	uqc := controllers.UserQueryController{QueryGateway: gf.GetQueryGateway()}
	ucc := controllers.UserCommandController{
		QueryGateway:   gf.GetQueryGateway(),
		CommandGateway: gf.GetCommandGateway(),
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Post("/users", ucc.HandleCreate)
	app.Get("/users/:id", uqc.HandleGetUserByIDQuery)

	log.Fatal(app.Listen(":8080"))
}
