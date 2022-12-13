package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
				Distance struct {
					Text string `json:"text"`
				} `json:"distance"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
	Status string `json:"status"`
}

type Maps2 struct {
	Routes []struct {
		Bounds struct {
			Northeast struct {
				Lat json.Number
				Lng json.Number
			} `json:"northeast" `
		} `json:"bounds"`
	} `json:"routes"`
	Status string `json:"status"`
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

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/maps", func(c *fiber.Ctx) error {
		APIKEY := "AIzaSyAatMhbCcoDEvvfUp89QqDWD96RkxmVydA"
		//fmt.Println(APIKEY)
		origin := "Ptuj"
		destination := "Maribor"
		params := "&units=metrics&avoidTolls=True&mode=walking"
		url := "https://maps.googleapis.com/maps/api/directions/json?origin=" + origin + "&destination=" + destination + params + "&key=" + APIKEY

		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
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
		if err := json.Unmarshal(body, &mapa); err != nil { // Parse []byte to go struct pointer
			fmt.Println(err)
			fmt.Println("Can not unmarshal JSON")
		}
		Koordinata_Lat := mapa.Routes[0].Legs[0].StartLocation.Lat
		Koordinata_Lng := mapa.Routes[0].Legs[0].StartLocation.Lng
		//print Lat and Lng
		//fmt.Println(Koordinata_Lat)
		//fmt.Println(Koordinata_Lng)
		vrni := "zacetek:" + origin + "Lat " + strconv.FormatFloat(Koordinata_Lng, 'E', -1, 64) + "Lng: " + strconv.FormatFloat(Koordinata_Lat, 'E', -1, 64) + " destinacija:" + destination
		return c.Send([]byte(vrni))
	})
	app.Get("/api/mapsDummy", func(c *fiber.Ctx) error {
		return c.SendString("koordinata je: " + string("69"))
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Maps api container working"))
	})

	app.Listen(":8002")
}
