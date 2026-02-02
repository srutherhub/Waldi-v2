package main

import (
	"fmt"
	"log"
	"net/http"
	"waldi-v2/controllers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()
	serveStaticFiles(mux)

	controllers.Init(mux)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func serveStaticFiles(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
