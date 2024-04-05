package routes

import (
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/views"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)


func render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

func SetRoutes(app *fiber.App) {
	app.Get("/", func ( c *fiber.Ctx) error {
		return render(c, views.Home())
	})

	app.Get("/login", func ( c *fiber.Ctx) error {
		return render(c, views.Login())
	})
}
