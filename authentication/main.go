package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Role     int    `json:"role"`
	Password string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"name"`
	Password string `json:"-"`
}

type arrayHealthCheck struct {
	Id     string         `json:"id"`
	Health []healthCheck2 `json:"types"`
}
type healthCheck2 struct {
	// Name of the health check
	Name string `json:"name"`
	// Status of the health check
	Status string `json:"status"`
	// Error message of the health check
	Error []string `json:"error"`
	// Timestamp of the health check
	Timestamp string `json:"timestamp"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func authentication() func(c *fiber.Ctx) error {

	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized)
		},
		SigningKey: []byte("secret"),
	})
}

func register(c *fiber.Ctx) error {
	req := new(User)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
	}
	var user User
	user.Name = req.Name
	user.Email = req.Email
	user.Username = req.Username
	if err := user.HashPassword(user.Password); err != nil {
		return err
	}

	//save this info in the database

	//create token and return it
	return nil
}

func main() {
	health := healthCheck2{
		Name:      "Connection",
		Status:    "No test",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	//TODO use .env variable
	var dsn string
	dsn = "postgres://zlqwvdmx:x0tl7AVnX4zi0rsqeKcf8R2dhjvqOpib@ella.db.elephantsql.com/zlqwvdmx"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		health.Status = "Error"
		health.Error = append(health.Error, err.Error())
		panic("failed to connect database")

	} else {
		health.Status = "Ok"
		health.Error = append(health.Error, "None")
	}
	db.AutoMigrate(User{})

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		healthC := healthCheck2{
			Name:      "Container",
			Status:    "OK",
			Error:     []string{"None"},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		healthAr := arrayHealthCheck{
			Id:     "Authentication",
			Health: []healthCheck2{healthC, health},
		}

		healt_json, err := json.Marshal(healthAr) // back to json

		if err != nil {
			panic(err)
		}
		return c.SendString(string(healt_json))
	})
	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}
		if req.Username == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
		}
		user := new(User)

		if err := db.Where("username = ?", req.Username).Find(user).Error; err != nil {
			return err
		}
		if err := user.CheckPassword(req.Password); err != nil {
			return err
		} else {
			//create token and return it
			// Create the Claims
			claims := jwt.MapClaims{
				"name": user.Username,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			}

			// Create token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			return c.JSON(fiber.Map{"token": t, "user": user})
		}
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		req := new(User)
		if err := c.BodyParser(req); err != nil {
			return err
		}
		if req.Username == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
		}
		var user User
		user.Name = req.Name
		user.Email = req.Email
		user.Username = req.Username
		user.Password = req.Password
		if err := user.HashPassword(user.Password); err != nil {
			return err
		}
		//save this info in the database
		db.Create(&user)

		//create token and return it
		// Create the Claims
		claims := jwt.MapClaims{
			"name": user.Username,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t, "user": user})
	})

	// Unauthenticated route
	//app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	// Restricted Routes
	app.Get("/authenticate", authentication(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})

	app.Listen(":8003")
}
