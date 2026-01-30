package main

import (
	"log"
	"os"

	"proundmhee/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := app.NewServer()

	log.Printf("proundmhee listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
