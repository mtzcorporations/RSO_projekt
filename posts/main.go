package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type Post struct {
	Id       uint      `json:"id"`
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Comments []Comment `json:"comments" gorm:"-" default:"[]"`
}

type Comment struct {
	Id     uint
	PostId string
	Text   string
}

func main() {
	//TODO use .env variable
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

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		var posts []Post
		db.Find(&posts)

		for i, post := range posts {
			response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/posts/%d/comments", post.Id))
			if err != nil {
				return err
			}

			var comments []Comment
			json.NewDecoder(response.Body).Decode(&comments)

			posts[i].Comments = comments
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
	app.Listen(":8000")
}
