package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var Info info

type info struct {
	Cluster string `json:"cluster"`
	Env     string `json:"env"`
	Region  string `json:"region"`
}

func setInfo() {
	c := os.Getenv("CLUSTER")
	e := os.Getenv("ENV")
	r := os.Getenv("REGION")

	Info = info{Cluster: c, Env: e, Region: r}
}

func healthz(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "OK"})
}

func getInfo(c *fiber.Ctx) error {
	return c.JSON(Info)
}

func main() {
	setInfo()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Rest API example",
	})

	app.Get("/healthz", healthz)

	// Middleware
	api := app.Group("/", logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	api.Get("/api/info", getInfo)

	log.Fatal(app.Listen(":8080"))
}
