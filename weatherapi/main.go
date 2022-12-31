package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type weather struct {
	Name string `json:"name"`
	Main struct {
		Kelvin   float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Oblaki []struct {
		Sonce string `json:"main"`
	} `json:"weather"`
}
type arrayHealthCheck struct {
	Id     string        `json:"id"`
	Health []healthCheck `json:"types"`
}
type healthCheck struct {
	// Name of the health check
	Name string `json:"name"`
	// Status of the health check
	Status string `json:"status"`
	// Error message of the health check
	Error []string `json:"error"`
	// Timestamp of the health check
	Timestamp string `json:"timestamp"`
}

func sendMetrics(timeElapsed string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsage := strconv.Itoa(int(m.Sys))
	base_url := "http://104.45.183.75/api/metrics/weather/"
	apiURL := base_url + timeElapsed[:len(timeElapsed)-2] + "/" + memoryUsage
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
}

func main() {
	health := healthCheck{
		Name:      "Api connection",
		Status:    "No test",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/weather", func(c *fiber.Ctx) error {
		start := time.Now()
		//url := "http://api.openweathermap.org/data/2.5/weather?lat=46.05&lon=14.50&units=metric&appid=ab8428d16bce2694fb18fbab32071873"
		url := "https://api.openweathermap.org/data/2.5/weather?lat=46.05&lon=14.50&appid=ab8428d16bce2694fb18fbab32071873"

		spaceClient := http.Client{
			Timeout: time.Second * 2, // Timeout after 2 seconds
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
		} else {
			health.Status = "OK"
			health.Error = []string{"None"}
			health.Timestamp = time.Now().Format(time.RFC3339)
		}

		req.Header.Set("User-Agent", "test")
		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
			health.Status = "ERROR"
			health.Error = append(health.Error, getErr.Error())
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
			health.Status = "ERROR"
			health.Error = append(health.Error, readErr.Error())
		}
		weather_lj := weather{}
		jsonErr := json.Unmarshal([]byte(body), &weather_lj) // json to our "weather" struct
		if jsonErr != nil {
			log.Fatal(jsonErr)
			health.Status = "ERROR"
			health.Error = append(health.Error, jsonErr.Error())
		}
		//fmt.Println(weather_lj.Name, weather_lj.Main, weather_lj.Oblaki[0])

		// send to metrics
		timeElapsed := time.Since(start).String()
		sendMetrics(timeElapsed)

		if weather_lj.Main.Humidity > 50 {
			return c.Send([]byte("Avto"))
		} else {
			return c.Send([]byte("Kolo"))
		}

	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Weather api container working"))
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		healthC := healthCheck{
			Name:      "Container",
			Status:    "OK",
			Error:     []string{"None"},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		healthAr := arrayHealthCheck{
			Id:     "weatherapi",
			Health: []healthCheck{healthC, health},
		}

		healt_json, err := json.Marshal(healthAr) // back to json
		if err != nil {
			panic(err)
		}
		return c.SendString(string(healt_json))
	})
	app.Listen(":8001")

}
