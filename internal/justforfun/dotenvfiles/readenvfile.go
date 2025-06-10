package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func basicLoad() {
	// By default, godotenv, search .env file at the root of a project
	dir, _ := os.Getwd()
	fmt.Printf("Searching .env file in: %s\n", dir)

	err := godotenv.Load("./internal/justforfun/dotenvfiles/.env")
	if err != nil {
		log.Printf("Error loading .env: %v\n", err)
		return
	}

	apiKey := os.Getenv("API_KEY")
	dbHost := os.Getenv("DB_HOST")

	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Printf("DB Host: %s\n", dbHost)
}

func loadAndOverwrite() {
	// By default, godotenv, search .env file at the root of a project
	dir, _ := os.Getwd()
	fmt.Printf("Searching .env file in: %s\n", dir)

	// Original value
	fmt.Printf("Current value: %s\n", os.Getenv("API_KEY"))

	// Overwrite current with .env value
	err := godotenv.Overload("./internal/justforfun/dotenvfiles/.env.new")
	if err != nil {
		log.Printf("Error loading .env: %v\n", err)
		return
	}

	// New value after load .env
	fmt.Printf("New value: %s\n", os.Getenv("API_KEY"))
}

func main() {
	basicLoad()
	loadAndOverwrite()
}
