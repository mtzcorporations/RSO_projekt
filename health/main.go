package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"net/http"
	"time"
)

type healthCheck struct {
	Id     string `json:"id"`
	Health []struct {
		Name string `json:"name"`
		// Status of the health check
		Status string `json:"status"`
		// Error message of the health check
		Error []string `json:"error"`
		// Timestamp of the health check
		Timestamp string `json:"timestamp"`
	} `json:"types"`
	// Name of the health check

}

//	func test(url string, txt string) string {
//		heal := healthCheck{
//			Name:      "Container",
//			Status:    "ERROR",
//			Error:     nil,
//
//		}
//		spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
//		req, err := http.NewRequest(http.MethodGet, url, nil)
//		if err != nil {
//			heal.Error = append(heal.Error, txt + err.Error())
//			heal.Timestamp = time.Now().Format(time.RFC3339)
//			// convert heal to string
//			healString, _ := json.Marshal(heal)
//			return string(healString)
//		}
//		req.Header.Set("User-Agent", "test")
//		res, getErr := spaceClient.Do(req)
//		if getErr != nil {
//			heal.Error = append(heal.Error, txt + getErr.Error())
//			heal.Timestamp = time.Now().Format(time.RFC3339)
//			// convert heal to string
//			healString, _ := json.Marshal(heal)
//			return string(healString)
//		}
//		if res.Body != nil {
//			defer res.Body.Close()
//		}
//
//		body, readErr := ioutil.ReadAll(res.Body)
//		if readErr != nil {
//			heal.Error = append(heal.Error, txt + readErr.Error())
//			heal.Timestamp = time.Now().Format(time.RFC3339)
//			// convert heal to string
//			healString, _ := json.Marshal(heal)
//			return string(healString)
//		}
//		return string(body)
//	}
func checkFNC(url string, txt string) int {
	spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 500
	}
	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {

		return 500
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return 500
	}
	//body to arrayHealthCheck
	var array healthCheck
	json.Unmarshal(body, &array)
	if array.Health[0].Status != "OK" {
		return 500
	} else {
		return 200
	}
}
func checkAut() string {
	url := "http://104.45.183.75/authentication"
	//url := "http://authentication:8003/"                  //weather
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
	url := "http://104.45.183.75/api/maps/health"
	//url := "http://mapsapi:8002/health"                   //weather
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
	url := "http://104.45.183.75/api/gas/health"
	//url := "http://gasapi:8004/health"                    //weather
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
	url := "http://104.45.183.75/api/weather/health"
	// url := "http://weatherapi:8001/health"                //weather
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
		//return c.SendString(weatherHealth)

	})
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Health Check - api container working")
		return c.SendString("Dela")

	})
	//started := time.Now()
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		res := checkFNC("http://104.45.183.75/authentication", "Authentication")
		w.WriteHeader(200)
		if res == 200 {
			w.Write([]byte("ok ; time: " + time.Now().String()))
		} else {
			w.Write([]byte("error; time: " + time.Now().String()))
		}
		//duration := time.Now().Sub(started)
		//if duration.Seconds() > 10 {
		//	w.WriteHeader(500)
		//	w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
		//} else {
		//	w.WriteHeader(200)
		//	w.Write([]byte("ok"))
		//}
		//started = time.Now()
	})
	app.Listen(":8080")

}
