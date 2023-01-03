package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strconv"
)

type Metrics struct {
	serviceName      string  `json:"name"`
	numApiCalls      int     `json:"numapicalls"`
	totalTime        float64 `json:"totaltime"`
	totalMemoryUsage int     `json:"totalmemoryusage"`
}

// Initialize service metricsapi
var servicesMetrics = map[string]Metrics{
	"weather": {
		serviceName:      "weather service",
		numApiCalls:      0,
		totalTime:        0.0, // seconds
		totalMemoryUsage: 0},  // memory in kilobytes
	"maps": {
		serviceName:      "maps service",
		numApiCalls:      0,
		totalTime:        0.0,
		totalMemoryUsage: 0},
	"posts": {
		serviceName:      "posts service",
		numApiCalls:      0,
		totalTime:        0.0,
		totalMemoryUsage: 0},
	"gas": {
		serviceName:      "gas service",
		numApiCalls:      0,
		totalTime:        0.0,
		totalMemoryUsage: 0},
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
		if servicesMetrics[service_name].numApiCalls > 0 {
			numOfApiCallsStr := "Št klicev na " + service_name + " je: " + strconv.Itoa(servicesMetrics[service_name].numApiCalls) + "\n"
			averageTimeStr := "Povprečen čas " + service_name + " je: " + fmt.Sprintf("%f", servicesMetrics[service_name].totalTime/float64(servicesMetrics[service_name].numApiCalls)) + " ms\n"
			averageMemUsgStr := "Povprečna poraba spomina " + service_name + " je: " + strconv.Itoa(servicesMetrics[service_name].totalMemoryUsage/servicesMetrics[service_name].numApiCalls) + " MB\n"
			return c.Send([]byte(numOfApiCallsStr + averageTimeStr + averageMemUsgStr))
		} else {
			return c.Send([]byte("Na mikrostoritvi " + servicesMetrics[service_name].serviceName + " še ni bilo aktivnosti"))
		}
	})
	app.Post("/:service_name/:time_elapsed/:memmory_usage", func(c *fiber.Ctx) error {
		service_name := c.Params("service_name")
		time_elapsed := c.Params("time_elapsed")
		memmory_usage := c.Params("memmory_usage")

		currServ := servicesMetrics[service_name]
		if timeEl, err := strconv.ParseFloat(time_elapsed, 64); err == nil {
			currServ.totalTime += timeEl
		}
		if memeUs, err := strconv.Atoi(memmory_usage); err == nil {
			if memeUs > 0 {
				currServ.totalMemoryUsage += memeUs / 1000000
			}
		}
		currServ.numApiCalls += 1
		servicesMetrics[service_name] = currServ

		return c.Send([]byte("Sucess"))
	})
	app.Get("/healthL", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":8005")

}
