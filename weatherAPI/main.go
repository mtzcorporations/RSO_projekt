package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type weather struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Oblaki []struct {
		Sonce string `json:"main"`
	} `json:"weather"`
}

func main() {
	url := "http://api.openweathermap.org/data/2.5/weather?lat=46.05&lon=14.50&units=metric&appid=ab8428d16bce2694fb18fbab32071873"
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	weather_lj := weather{}
	jsonErr := json.Unmarshal([]byte(body), &weather_lj) // json to our "weather" struct
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	//fmt.Println(weather_lj.Name, weather_lj.Main, weather_lj.Oblaki[0])
	weather_lj_json, jsonErr := json.Marshal(weather_lj) // back to json
	if err != nil {
		log.Fatal(jsonErr)
	}

	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/test", func(c *fiber.Ctx) error {
		return c.SendString(string(weather_lj_json))
	})

	app.Listen(":8001")

}
