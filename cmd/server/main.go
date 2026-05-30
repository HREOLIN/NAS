package main

import (
	"log"

	"github.com/HREOLIN/NAS/internal/app"
)

func main() {
	server := app.NewServer(":8080")
	log.Println("NAS backend demo is running at http://127.0.0.1:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
