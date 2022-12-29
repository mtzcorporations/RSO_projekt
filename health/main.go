package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"net/http"
	"time"
)

func checkAut() string {
	//url := "http://10.0.143.93:8004/health"
	url := "http://authentication:8003/"                  //weather
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return "Authentication: " + err.Error()
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Println(err.Error())
		return "Authentication " + getErr.Error()
	}
	if res.Body != nil {

		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		fmt.Println(err.Error())
		return "Authentication: " + readErr.Error()
	}
	return string(body)
}
func checkMaps() string {
	//url := "http://10.0.143.93:8004/health"
	url := "http://mapsapi:8002/health"                   //weather
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "MapsApi: " + err.Error()
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return "MapsApi: " + getErr.Error()
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "MapsApi: " + readErr.Error()
	}
	return string(body)
}
func checkGas() string {
	//url := "http://10.0.143.93:8004/health"
	url := "http://gasapi:8004/health"                    //weather
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "GasApi: " + err.Error()
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return "GasApi: " + getErr.Error()
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "GasApi: " + readErr.Error()
	}
	return string(body)
}
func checkWeather() string {
	//url := "http://10.0.143.93:8001/health"
	url := "http://weatherapi:8001/health"                //weather
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "WeatherApi: " + err.Error()
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return "WeatherApi: " + getErr.Error()
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "WeatherApi: " + readErr.Error()
	}
	return string(body)
}
func main() {

	// Expose API
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		weatherHealth := checkWeather()
		GasHealth := checkGas()
		MapsHealth := checkMaps()
		AuthHealth := checkAut()
		return c.SendString(weatherHealth + "\n " + GasHealth + "\n " + MapsHealth + "\n " + AuthHealth)

	})

	app.Listen(":8080")

}
