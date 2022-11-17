package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Post struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type Comment struct {
	Id     uint
	PostId string
	Text   string
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
		url := "http://weatherapi:8001/"
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
		return c.SendString(string(body))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Posts container working"))
	})

	app.Listen(":8000")
}
