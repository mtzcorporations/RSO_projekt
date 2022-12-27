package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strconv"
)

type metrics struct {
	serviceName string `json:"name"`
	numApiCalls int    `json:"numapicalls"`
}

func main() {

	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	// Initialize service metrics
	var servicesMetrics = map[string]*metrics{
		"weather": {
			serviceName: "weather",
			numApiCalls: 0},
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Metrics container working"))
	})

	// GET/POST requests

	app.Get("/:service_name/calls", func(c *fiber.Ctx) error {
		service_name := c.Params("service_name")
		return c.SendString("Št klicev na " + service_name + " service je: " + strconv.Itoa(servicesMetrics["weather"].numApiCalls))
	})
	app.Post("/:service_name/calls", func(c *fiber.Ctx) error {
		servicesMetrics["weather"].numApiCalls += 1
		return c.Send([]byte("Sucess"))
	})

	app.Listen(":8005")

}
