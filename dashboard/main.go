package main

import (
	"database/sql"
	graphdb "optitraffic/graphDB"
	"optitraffic/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "../graphDB/db.db")
	if err != nil {
		log.Fatal(err)
	}
	dao := graphdb.NewDAO(db)
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

	app.Get("/points", func(c *fiber.Ctx) error {
		points, err := dao.GetAllPoints()
		s, err := points.ToJSON()
		if err != nil {
			return err
		}

		return c.Send([]byte(s))
	})

	app.Get("/lines", func(c *fiber.Ctx) error {
		lines, err := dao.GetAllPaths()
		s, err := lines.ToJSON()
		if err != nil {
			return err
		}

		return c.Send([]byte(s))
	})

	app.Get("/vehicles", func(c *fiber.Ctx) error {
		return nil
	})

	app.Listen(":6969")
}
