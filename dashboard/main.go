package main

import (
	"database/sql"
	graphconvertor "optitraffic/graphConvertor"
	graphdb "optitraffic/graphDB"
	"optitraffic/node"
	"optitraffic/templates"
	"optitraffic/traffic"

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

	graph, err := dao.GetGraph()
	if err != nil {
		log.Fatal(err)
	}

	tm := traffic.NewTrafficManager(&graph)

	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)
	tm.NewRandomVehicle(node.EmergencyVehicle)

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
		tm.Update(0.0001)
		fc := tm.VehiclesAsPoints()
		s, err := fc.ToJSON()
		if err != nil {
			return err
		}

		return c.Send([]byte(s))
	})

	app.Get("/lights", func(c *fiber.Ctx) error {
		path, _ := graphconvertor.TurnGraphToGeoJSON(*tm.Graph)
		path_s, err := path.ToJSON()
		if err != nil {
			return err
		}
		return c.Send([]byte(path_s))
	})

	app.Listen(":6969")
}
