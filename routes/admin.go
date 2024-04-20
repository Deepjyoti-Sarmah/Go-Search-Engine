package routes

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/db"
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/utils"
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/views"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Something went wrong</h2>")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2> invalid password</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2> Error: Something went wrong logging in</h2>")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")

	return c.SendStatus(200)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	//Get the cookie by name
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}
	// Parse the cookie & check for errors
	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", 302)
	}
	//Parse the custom claims & check jwt is valid
	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/login", 302)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := &db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type settingForm struct {
	Amount   int  `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func DashboardPostHandler ( c *fiber.Ctx) error {
  input := settingForm{}
  if err := c.BodyParser(&input); err != nil {
	c.Status(500)
	return c.SendString("<h2>Error: Something went wrong</h2>")
  }

  addNew := false
  if input.AddNew == "on" {
	addNew = true
  }

  searchOn := false
  if input.SearchOn == "on" {
	searchOn = true
  }

  settinngs := &db.SearchSetting{}
  settinngs.Amount = uint(input.Amount)
  settinngs.SearchOn = searchOn
  settinngs.AddNew = addNew
  err := settinngs.Update()
  if err != nil {
	fmt.Println(err)
	return c.SendString("<h2>Error: Cannot update settings</h2>")
  }

  c.Append("HX-Refresh", "true")
  return c.SendStatus(200)
}

