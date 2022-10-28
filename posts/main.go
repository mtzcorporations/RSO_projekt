package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func main() {
	//TODO use .env variable
	//dsn := "root@tcp(127.0.0.1:3306)/posts_ms"
	//dsn := "root:root@tcp(postsmysql:3306)/post_ms"
	time.Sleep(20 * time.Second)
	dsn := "tester:secret@tcp(postsmysql:3306)/test"
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
