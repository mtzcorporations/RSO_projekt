package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

// @Summary User object
// @Description Represents a user in the system
// @Tags User
type User struct {
	// The unique ID of the user
	Id uint `json:"id"`
	// The user's name
	Name string `json:"name"`
	// The user's username
	Username string `json:"username" gorm:"unique"`
	// The user's email
	Email string `json:"email" gorm:"unique"`
	// The user's role
	Role int `json:"role"`
	// The user's password (hashed)
	Password string `json:"password"`
}

// @Summary Login request object
// @Description Represents a request to log in to the system
// @Tags Login
type LoginRequest struct {
	// The user's username
	Username string `json:"username"`
	// The user's password
	Password string `json:"password"`
}

// @Summary Health check object
// @Description Represents the status of a health check
// @Tags Health Check
type arrayHealthCheck struct {
	// The unique ID of the health check
	Id string `json:"id"`
	// The list of health check types
	Health []healthCheck2 `json:"types"`
}

// @Summary Health check type object
// @Description Represents a type of health check
// @Tags Health Check
type healthCheck2 struct {
	// The name of the health check
	Name string `json:"name"`
	// The status of the health check
	Status string `json:"status"`
	// The error message of the health check
	Error []string `json:"error"`
	// The timestamp of the health check
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
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
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
		SigningKey: []byte(os.Getenv("JWT_KEY")),
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

	// @Summary Returns the health status of the service
	// @Tags Health
	// @Produce  json
	// @Success 200 {object} HealthResponse
	// @Router / [get]
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
	// @Summary Logs in a user
	// @Tags Login
	// @Accept  json
	// @Produce  json
	// @Param request body object true "LoginRequest"
	// @Success 200 {object} LoginResponse
	// @Failure 400 {string} string "Invalid credentials"
	// @Failure 500 {string} string "Internal server error"
	// @Router /login [post]
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
			t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			return c.JSON(fiber.Map{"token": t, "user": user})
		}
	})

	// @Summary Registers a new user
	// @Tags Registration
	// @Accept  json
	// @Produce  json
	// @Param request body object true "User"
	// @Success 200 {object} RegisterResponse
	// @Failure 400 {string} string "Invalid credentials"
	// @Failure 500 {string} string "Internal server error"
	// @Router /register [post]
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
		//get token name

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t, "user": user})
	})

	// Unauthenticated route
	//app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}))

	// Restricted Routes
	// @Summary Verifies the user's token
	// @Tags Authentication
	// @Security JWT
	// @Success 200 {string} string "Accepted"
	// @Failure 401 {string} string "Unauthorized"
	// @Router /authenticate [get]
	app.Get("/authenticate", authentication(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
	//http.HandleFunc("/healthR", func(w http.ResponseWriter, r *http.Request) {
	//	if err != nil {
	//		w.WriteHeader(200)
	//		w.Write([]byte("ok ; time: " + time.Now().String()))
	//	} else {
	//		w.WriteHeader(500)
	//		w.Write([]byte("error; time: " + time.Now().String()))
	//	}
	//})

	// @Summary Checks the server's status
	// @Tags Health Check
	// @Success 200 {string} string "OK"
	// @Failure 500 {string} string "Internal server error"
	// @Router /healthR [get]
	app.Get("/healthR", func(c *fiber.Ctx) error {
		fmt.Println(err)
		if err != nil {
			return c.SendStatus(500)
		} else {
			return c.SendStatus(200)
		}
	})

	// @Summary Checks the server's status
	// @Tags Health Check
	// @Success 200 {string} string "OK"
	// @Router /healthL [get]
	app.Get("/healthL", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":8003")
}
