package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/db"
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/routes"
	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("cannot find environmental variable")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())
	db.InitDB()

	routes.SetRoutes(app)
	utils.StartCornJobs()

	// Start our server and listen for a shutudown
	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block the main thread untill interrupted

	app.Shutdown()
	fmt.Println("shutting down server")
}
