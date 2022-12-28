package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func checkWeather() string {
	//url := "http://10.0.143.93:8001/health"
	url := "http://weatherapi:8001/health"                //weather
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return getErr.Error()
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return string(body)
}
func main() {

	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		weatherHealth := checkWeather()

		return c.SendString(weatherHealth)
	})

	app.Listen(":8080")

}
