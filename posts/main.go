package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	//TODO use .env variable

	// var dsn string
	// if true {
	// 	time.Sleep(5 * time.Second)
	// 	dsn = "tester:secret@tcp(postsmysql:3306)/test"
	// } else {
	// 	dsn = "root@tcp(127.0.0.1:3306)/posts_ms"
	// }
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// db.AutoMigrate(Post{})

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/test", func(c *fiber.Ctx) error {
		// Do api request to another container
		// url := "http://weatherapi:8001/api/test"
		url := "http://10.0.41.147:8001/api/test"
		spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
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

		degrees := strconv.FormatFloat(weather_lj.Main.Kelvin, 'E', -1, 64)
		weather_str := "V " + weather_lj.Name + " je " + degrees + " stopinj"
		return c.SendString(weather_str)
	})

	app.Get("/api/maps", func(c *fiber.Ctx) error {
		// Do api request to another container
		//url := "http://mapsapi:8002/api/mapsDummy"
		url := "http://10.0.182.147:8002/api/mapsDummy"
		spaceClient := http.Client{Timeout: time.Second * 20} // Timeout after 2 seconds
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

		return c.SendString("koordinata je: " + string(body))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Posts container working"))
	})

	app.Listen(":8000")
}
