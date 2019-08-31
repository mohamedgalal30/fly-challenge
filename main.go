package main

import (
	"log"
	"net/http"
	"os"

	"fly/handlers"
	_ "fly/providers"
)

func main() {
	PORT := getPort()

	handlers.Routes()

	log.Println("listener : Started : Listening on :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func getPort() string {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3000"
	}
	return PORT
}
