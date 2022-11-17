package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	/*
		var dsn string
		if os.Getenv("ISDOCKER") == "1" {
			time.Sleep(5 * time.Second)
			dsn = "tester:secret@tcp(postsmysql:3306)/test"
		} else {
			dsn = "root@tcp(127.0.0.1:3306)/posts_ms"
		}
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(Post{})
	*/

	app := fiber.New()

	app.Use(cors.New())

	/*
		app.Get("/api/posts", func(c *fiber.Ctx) error {
			var posts []Post
			db.Find(&posts)

			for i, post := range posts {
				response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/posts/%d/comments", post.Id))
				if err != nil {
					return err
				}

				var weatherapi []Comment
				json.NewDecoder(response.Body).Decode(&weatherapi)

				posts[i].Comments = weatherapi
			}

			return c.JSON(posts)
		})

		app.Post("/api/posts", func(c *fiber.Ctx) error {
			var post Post
			if err := c.BodyParser(&post); err != nil {
				return err
			}

			db.Create(&post)
			return c.JSON(post)
		})
	*/

	app.Get("/", func(c *fiber.Ctx) error {
		// Do api request to another container
		url := "http://localhost:8001/"
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

	app.Listen(":8000")
}
