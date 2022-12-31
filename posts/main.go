package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Rating   string `json:"rating"`
	Password string `json:"-"`
}
type RatingRequest struct {
	Id     uint   `json:"id"`
	Rating string `json:"rating"`
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

	var dsn string
	dsn = "postgres://zlqwvdmx:x0tl7AVnX4zi0rsqeKcf8R2dhjvqOpib@ella.db.elephantsql.com/zlqwvdmx"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Post("/rate/driver", func(c *fiber.Ctx) error {
		// Do api request to another container
		// url := "http://weatherapi:8001/api/test"
		req := new(User)
		if err := c.BodyParser(req); err != nil {
			return err
		}
		if req.Rating == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
		}
		//save this info in the database
		db.Model(&User{}).Where("id = ?", req.Id).Update("rating", req.Rating)
		return c.SendStatus(fiber.StatusAccepted)
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		result := []User{}
		db.Find(&result)
		return c.Status(http.StatusOK).JSON(&result)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON([]byte("Posts container working"))
	})

	app.Listen(":8000")
}
