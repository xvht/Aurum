package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"vexal/database"
	"vexal/handlers"
	"vexal/lib"
	"vexal/packages/gitinfo"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	commit, err := gitinfo.GetCommit()
	if err != nil {
		log.Fatalf("Error getting commit: %v", err)
	}
	log.Printf("Commit: %s", commit)

	log.Printf("--- Booting Vexal Key Server --- ")

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbUrl := os.Getenv("DATABASE_URL")
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	log.Printf("Database: %s", lib.GetDisplayURL(dbUrl))
	database.Connect(dbUrl)

	app := fiber.New()
	for path, handler := range handlers.Handlers {
		parts := strings.Split(strings.TrimSpace(path), " ")
		method := parts[0]
		route := parts[1]

		switch method {
		case "GET":
			app.Get(route, handler)
		case "POST":
			app.Post(route, handler)
		case "PUT":
			app.Put(route, handler)
		case "DELETE":
			app.Delete(route, handler)
		case "PATCH":
			app.Patch(route, handler)
		case "OPTIONS":
			app.Options(route, handler)
		case "HEAD":
			app.Head(route, handler)
		default:
			log.Fatalf("Invalid method: %s", method)
		}

		log.Printf("Router: %s %s registered", method, route)
	}

	log.Printf("Server started on %s", port)
	log.Fatal(app.Listen(port, fiber.ListenConfig{DisableStartupMessage: true}))
}
