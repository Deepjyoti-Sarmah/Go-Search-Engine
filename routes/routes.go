package routes

import (
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/db"
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

type settingForm struct {
	Amount   int  `form:"amount"`
	SearchOn bool `form:"searchOn"`
	AddNew   bool `form:"addNew"`
}

func SetRoutes(app *fiber.App) {
	app.Get("/",AuthMiddleware, LoginHandler)
	app.Post("/", AuthMiddleware, LoginPostHandler)

	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)

  app.Get("/create", func(c *fiber.Ctx) error {
	u := &db.User{}
	u.CreateAdmin()
	return c.SendString("created")
  })
}
