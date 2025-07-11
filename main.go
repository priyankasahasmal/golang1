package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	HOST     = "dpg-d1l9vbbe5dus73feni1g-a.oregon-postgres.render.com"
	PORT     = "5432"
	USERNAME = "database1_mi2h_user"
	PASSWORD = "jqy7Bg3KKql230L8LMLJfArda5ziDWPS"
	DBNAME   = "database1_mi2h"
)

func GetPsqlInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		HOST, PORT, USERNAME, PASSWORD, DBNAME)
}

func CreateDbObject() error {
	var err error
	DB, err = sql.Open("postgres", GetPsqlInfo())
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetConnMaxLifetime(1 * time.Hour)

	return nil
}

type FetchAllUsersOutput struct {
	UserID int64
	Name   string
	Email  string
}

func FetchAllUsers() ([]FetchAllUsersOutput, error) {
	query := "SELECT userid, name, email FROM users"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []FetchAllUsersOutput
	for rows.Next() {
		var u FetchAllUsersOutput
		err := rows.Scan(&u.UserID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func main() {
	if err := CreateDbObject(); err != nil {
		log.Fatal("❌ DB Connection Failed:", err)
	}

	app := fiber.New()

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Go app is successfully deployed on Render!")
	})

	// Users list
	app.Get("/users", func(c *fiber.Ctx) error {
		users, err := FetchAllUsers()
		if err != nil {
			return c.Status(500).SendString("❌ Failed to fetch users: " + err.Error())
		}
		return c.JSON(users)
	})

	// Sign-in route
	app.Post("/signin", SignInWeb)

	// ✅ POST /get-userid route
	app.Post("/get-userid", func(c *fiber.Ctx) error {
		var req struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString("❌ Invalid JSON")
		}

		userID, err := FetchUserIDFromEmailID(req.Email)
		if err != nil {
			return c.Status(500).SendString("❌ " + err.Error())
		}

		return c.JSON(fiber.Map{
			"email":  req.Email,
			"userid": userID,
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInWeb(c *fiber.Ctx) error {
	signInRequestObject := SignInRequest{}

	// Parse request body
	if err := c.BodyParser(&signInRequestObject); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Fetch user ID from DB
	userID, err := FetchUserIDFromEmailID(signInRequestObject.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "User found",
		"userId":  userID,
		"email":   signInRequestObject.Email,
	})
}

func FetchUserIDFromEmailID(email string) (int, error) {
	query := `SELECT userid FROM users WHERE email = $1`
	row := DB.QueryRow(query, email)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
