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

	r, deps, err := app.NewServer(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatal(err)
	}
	defer deps.Log.Sync()

	log.Printf("proundmhee listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
