package main

import (
	"database/sql"
	"fmt"
	"optitraffic/geojson"
	graphconvertor "optitraffic/graphConvertor"
	graphdb "optitraffic/graphDB"
	"optitraffic/templates"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	_ "github.com/mattn/go-sqlite3"
)

var creating_graph = false

func geojson_string_pair_parse(input string) [2]string {
	// Split the input string at the boundary between two FeatureCollections
	parts := strings.Split(input, `},{"type":"FeatureCollection"`)

	// Add back the missing curly braces around the two parts
	part1 := parts[0] + `}`
	part1 = part1[1:]

	part2 := `{"type":"FeatureCollection"` + parts[1]
	part2 = part2[:len(part2)-1]

	// Create an array with the two strings
	return [2]string{part1, part2}
}

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
		parts := geojson_string_pair_parse(string(c.Body()))

		resp_points, err := geojson.FeatureCollFromJSON[geojson.FlatGeometry](parts[0])
		if err != nil {
			return err
		}

		resp_lines, err := geojson.FeatureCollFromJSON[geojson.Geometry](parts[1])
		if err != nil {
			return err
		}

		resp_nodes, err := graphconvertor.PointsCollToGeoNode(resp_points)
		if err != nil {
			return err
		}
		err = dao.StoreGeoNodes(resp_nodes...)
		if err != nil {
			fmt.Printf("shit")
			return err
		}

		resp_paths, err := graphconvertor.LineCollToGeoPath(resp_lines)
		if err != nil {
			return err
		}
		err = dao.StoreGeoPaths(resp_paths...)
		if err != nil {
			fmt.Printf("fuck")
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":6969")
}
