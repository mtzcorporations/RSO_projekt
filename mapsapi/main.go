package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Maps struct {
	Routes []struct {
		Legs []struct {
			Distance struct {
				Text string `json:"text"`
			} `json:"distance"`
			Duration struct {
				Text string `json:"text"`
			} `json:"duration"`
			StartLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"start_location"`
			Steps []struct {
				EndLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"end_location"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
	Status string `json:"status"`
}

type Mapsout struct {
	Zacetek    string `json:"zacetek"`
	Konec      string `json:"cilj"`
	Trajanje   string `json:"trajanje"`
	Razdalja   string `json:"razdalja"`
	Koordinate []struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
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

type Path struct {
	Path string `json:"path"`
}

// string function , returning string
func tipiPoti(pot string) string {
	// če je pot Peš vrni walking, če je pot z avtomobilom vrni driving, če je pot z vlakom vrni transit, če je kolo vrni bicycling
	if pot == "Peš" {
		return "walking"
	} else if pot == "Vlak" {
		return "transit"
	} else if pot == "Kolo" {
		return "bicycling"
	}
	return "driving"
}

func sendMetrics(timeElapsed string) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsage := strconv.Itoa(int(m.Sys))
	base_url := "http://104.45.183.75/api/metrics/maps/"
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
		return
	}
	defer res.Body.Close()
}

func main() {

	health := healthCheck{
		Name:      "Api connection",
		Status:    "No test",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	cb := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "my-circuit-breaker",
			MaxRequests: 2,
			Timeout:     2 * time.Second,
			Interval:    20 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				fmt.Println(counts.ConsecutiveFailures > 2)
				return counts.ConsecutiveFailures > 2
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				fmt.Printf("CircuitBreaker '%s' changed from '%s' to '%s'\n", name, from, to)
			},
		},
	)

	//getApiDat_testFunc()
	app := fiber.New()

	app.Use(cors.New())

	app.Post("/test", func(c *fiber.Ctx) error {

		start := time.Now()

		path_req := new(Path)
		if err := c.BodyParser(path_req); err != nil {
			return err
		}
		locations := path_req.Path

		if locations == "" {
			return nil
		}

		locations_Array := strings.Split(locations, "|")

		origin := locations_Array[0]
		destination := locations_Array[len(locations_Array)-1]

		locations_between := ""
		for i := 1; i < len(locations_Array)-1; i++ {
			locations_between = locations_between + locations_Array[i]
			if len(locations_Array)-2 != i {
				locations_between = locations_between + "|"
			}
		}

		APIKEY := os.Getenv("API_KEY")
		//origin := "Ptuj"
		fmt.Println(locations_between)
		waypoints := "&waypoints=" + locations_between // | je ločilo med waypointi
		//destination := "Maribor"
		params := "&units=metricsapi&mode=driving" // TODO WARNING maybe wrong refactor
		apiUrl := "https://maps.googleapis.com/maps/api/directions/json?origin=" + origin + "&destination=" + destination + waypoints + params + "&key=" + APIKEY
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, apiUrl, nil)
		if err != nil {
			fmt.Println("empty")
			fmt.Println(err)
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
			return err

		} else {
			health.Status = "OK"
			health.Error = []string{"None"}
			health.Timestamp = time.Now().Format(time.RFC3339)
		}
		body, err := cb.Execute(func() (interface{}, error) {
			res, err := client.Do(req)
			if err != nil {
				fmt.Println("http Get request gave error")
				fmt.Println(err)
				return nil, err
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				health.Status = "ERROR"
				health.Error = append(health.Error, err.Error())
				health.Timestamp = time.Now().Format(time.RFC3339)
				return nil, err
			}
			return body, nil

		})

		if err != nil {
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
			return err
		}

		//desifriranje jsona
		var mapa Maps
		var output Mapsout
		if err := json.Unmarshal(body.([]byte), &mapa); err != nil { // Parse []byte to go struct pointer
			fmt.Println(err)
			fmt.Println("Can not unmarshal JSON")
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
			return err
		}

		output.Razdalja = mapa.Routes[0].Legs[0].Distance.Text
		output.Trajanje = mapa.Routes[0].Legs[0].Duration.Text
		output.Zacetek = origin
		output.Konec = destination
		output.Koordinate = append(output.Koordinate, struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{mapa.Routes[0].Legs[0].StartLocation.Lat, mapa.Routes[0].Legs[0].StartLocation.Lng})
		//fmt.Println(len(mapa.Routes[0].Legs))
		for j := 0; j < len(mapa.Routes[0].Legs); j++ {
			for i := 0; i < len(mapa.Routes[0].Legs[j].Steps); i++ {
				output.Koordinate = append(output.Koordinate, struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				}{mapa.Routes[0].Legs[j].Steps[i].EndLocation.Lat, mapa.Routes[0].Legs[j].Steps[i].EndLocation.Lng})
			}
		}
		vrni, err := json.Marshal(output)
		if err != nil {
			fmt.Println(err)
			health.Status = "ERROR"
			health.Error = append(health.Error, err.Error())
			health.Timestamp = time.Now().Format(time.RFC3339)
		}

		// send to metricsapi
		timeElapsed := time.Since(start).String()
		sendMetrics(timeElapsed)

		// return
		return c.Send(vrni)
	})
	app.Get("/mapsDummy", func(c *fiber.Ctx) error {
		return c.SendString("koordinata je: " + string("69"))
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Maps api container working"))
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		healthC := healthCheck{
			Name:      "Container",
			Status:    "OK",
			Error:     []string{"None"},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		healthAr := arrayHealthCheck{
			Id:     "MapsApi",
			Health: []healthCheck{healthC, health},
		}

		healt_json, err := json.Marshal(healthAr) // back to json
		if err != nil {
			panic(err)
		}
		return c.SendString(string(healt_json))
	})

	app.Get("/healthL", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":8002")
}
