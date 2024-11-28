package main

import (
	"optitraffic/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var creating_graph = false

func main() {
	app := fiber.New()

	templates, err := templates.NewTemplates()
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/index.js", func(c *fiber.Ctx) error {
		return c.SendFile("./index.js")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return templates.Render("index", nil, c)
	})

	app.Get("/creator_not_creating", func(c *fiber.Ctx) error {
		return templates.Render("creator_not_creating", nil, c)
	})

	app.Get("/create_graph", func(c *fiber.Ctx) error {
		return templates.Render("creator_creating", nil, c)
	})

	app.Get("/save_graph", func(c *fiber.Ctx) error {
		return templates.Render("creator_not_creating", nil, c)
	})

	app.Post("/save_geojson", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":6969")
}
