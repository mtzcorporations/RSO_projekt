package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"reflect"
	"strconv"
)

type Metrics struct {
	serviceName string  `json:"name"`
	numApiCalls int     `json:"numapicalls"`
	totalTime   float64 `json:"totaltime"`
}

// Initialize service metrics
var servicesMetrics = map[string]Metrics{
	"weather": {
		serviceName: "weather service",
		numApiCalls: 0,
		totalTime:   0.0},
	"maps": {
		serviceName: "maps service",
		numApiCalls: 0,
		totalTime:   0.0},
	"posts": {
		serviceName: "posts service",
		numApiCalls: 0,
		totalTime:   0.0},
	"gas": {
		serviceName: "gas service",
		numApiCalls: 0,
		totalTime:   0.0},
}

func getFieldString(m *Metrics, field string) string {
	r := reflect.ValueOf(m)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func main() {

	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Metrics container working"))
	})

	// GET/POST requests
	app.Get("/:service_name", func(c *fiber.Ctx) error {
		service_name := c.Params("service_name")
		numOfApiCallsStr := "Št klicev na " + service_name + " je: " + strconv.Itoa(servicesMetrics[service_name].numApiCalls) + "\n"
		averageTimeStr := "Povprečen čas " + service_name + " je: " + fmt.Sprintf("%f", servicesMetrics[service_name].totalTime/float64(servicesMetrics[service_name].numApiCalls)) + "\n"
		return c.Send([]byte(numOfApiCallsStr + averageTimeStr))
	})
	app.Post("/:service_name/:time_elapsed", func(c *fiber.Ctx) error {
		service_name := c.Params("service_name")
		time_elapsed := c.Params("time_elapsed")

		currServ := servicesMetrics[service_name]
		if s, err := strconv.ParseFloat(time_elapsed, 64); err == nil {
			currServ.totalTime += s
		}
		currServ.numApiCalls += 1
		servicesMetrics[service_name] = currServ

		return c.Send([]byte("Sucess"))
	})

	app.Listen(":8005")

}
