package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type jsn struct {
	// body struct
	Result struct {
		Bencin string `json:"gasoline"`
		Dizel  string `json:"diesel"`
	} `json:"result"`
}
type jsnret struct {
	Bencin string `json:"bencin"`
	Dizel  string `json:"dizel"`
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

func getDataJson() {
	url := "https://api.collectapi.com/gasPrice/fromCity?city=ljubljana?currency=eur'"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "apikey 2M0y9SHvCFNV5KUD2lGZL2:3VnJ9JIwyF4UCf01Ffbx3S")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("empty")
		fmt.Println(err)

	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	var data jsn
	var retrn jsnret
	if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
		fmt.Println("Can not unmarshal JSON")
	}
	retrn.Dizel = data.Result.Dizel
	retrn.Bencin = data.Result.Bencin
	fmt.Println(retrn)
	vrni, err := json.Marshal(retrn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(vrni))
	//return vrni from function

}
func sendMetrics(timeElapsed string) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsage := strconv.Itoa(int(m.Sys))
	base_url := "http://104.45.183.75/api/metricsapi/gasapi/"
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
	//getDataJson()
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/gas", func(c *fiber.Ctx) error {
		start := time.Now()
		url := "https://api.collectapi.com/gasPrice/fromCity?city=ljubljana?currency=eur'"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "apikey 2M0y9SHvCFNV5KUD2lGZL2:3VnJ9JIwyF4UCf01Ffbx3S")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			//fmt.Println("empty")
			fmt.Println(err)
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
		} else {
			health.Status = "OK"
			health.Error = []string{"None"}
			health.Timestamp = time.Now().Format(time.RFC3339)
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())

		}
		var data jsn
		var retrn jsnret
		if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to go struct pointer
			fmt.Println(err)
			fmt.Println("Can not unmarshal JSON")
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
		}
		retrn.Dizel = data.Result.Dizel
		retrn.Bencin = data.Result.Bencin
		fmt.Println(retrn)
		vrni, err := json.Marshal(retrn)
		if err != nil {
			fmt.Println(err)
		}

		// send to metricsapi
		timeElapsed := time.Since(start).String()
		sendMetrics(timeElapsed)

		return c.Send(vrni)
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Gasoline api container working"))
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		healthC := healthCheck{
			Name:      "Container",
			Status:    "OK",
			Error:     []string{"None"},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		healthAr := arrayHealthCheck{
			Id:     "GasApi",
			Health: []healthCheck{healthC, health},
		}

		healt_json, err := json.Marshal(healthAr) // back to json
		if err != nil {
			panic(err)
		}
		return c.SendString(string(healt_json))
	})
	app.Listen(":8004")

}
