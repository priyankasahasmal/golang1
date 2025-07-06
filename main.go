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

	// Root route to confirm deployment
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Go app is successfully deployed on Render!")
	})

	// Endpoint to fetch users
	app.Get("/users", func(c *fiber.Ctx) error {
		users, err := FetchAllUsers()
		if err != nil {
			return c.Status(500).SendString("❌ Failed to fetch users: " + err.Error())
		}
		return c.JSON(users)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
