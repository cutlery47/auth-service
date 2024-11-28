package main

import (
	"log"

	"github.com/cutlery47/auth-service/internal/app"
)

func main() {
	log.Fatal("error: ", app.Run())
}
