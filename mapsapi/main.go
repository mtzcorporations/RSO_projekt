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

func fileReader2() string {
	content, err := ioutil.ReadFile("KEYS.TXT")

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
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
	// apiURL = "http://104.45.183.75/metrics/maps"
	apiURL := "http://metrics:8005/maps/" + timeElapsed[:len(timeElapsed)-2] + "/" + memoryUsage
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
	//getApiDat_testFunc()
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/test", func(c *fiber.Ctx) error {

		start := time.Now()
		// APIKEY := os.Getenv("API_KEY")
		APIKEY := "AIzaSyB8YSNqlWm6FMKuOfBnHL223E7m6Uate6Q"
		origin := "Ptuj"
		destination := "Maribor"
		params := "&units=metrics&avoidTolls=True&mode=driving"
		apiURL := "https://maps.googleapis.com/maps/api/directions/json?origin=" + origin + "&destination=" + destination + params + "&key=" + APIKEY
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, apiURL, nil)
		if err != nil {
			fmt.Println("empty")
			fmt.Println(err)

		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)

		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)

		}

		//desifriranje jsona
		var mapa Maps
		var output Mapsout
		if err := json.Unmarshal(body, &mapa); err != nil { // Parse []byte to go struct pointer
			fmt.Println(err)
			fmt.Println("Can not unmarshal JSON")
		}

		output.Razdalja = mapa.Routes[0].Legs[0].Distance.Text
		output.Trajanje = mapa.Routes[0].Legs[0].Duration.Text
		output.Zacetek = origin
		output.Konec = destination
		output.Koordinate = append(output.Koordinate, struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{mapa.Routes[0].Legs[0].StartLocation.Lat, mapa.Routes[0].Legs[0].StartLocation.Lng})

		for i := 1; i < len(mapa.Routes[0].Legs[0].Steps); i++ {
			output.Koordinate = append(output.Koordinate, struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			}{mapa.Routes[0].Legs[0].Steps[i].EndLocation.Lat, mapa.Routes[0].Legs[0].Steps[i].EndLocation.Lng})
		}
		vrni, err := json.Marshal(output)
		if err != nil {
			fmt.Println(err)
		}

		// send to metrics
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

	app.Listen(":8002")
}
