package templates

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(name string, data any, c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "text/html")
	return t.templates.ExecuteTemplate(c.Response().BodyWriter(), name, data)
}

func NewTemplates() (*Templates, error) {
	templates, err := template.ParseGlob("*.html")
	if err != nil {
		return nil, err
	}
	return &Templates{
		templates: templates,
	}, err
}
