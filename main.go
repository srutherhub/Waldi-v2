package main

import (
	"fmt"
	"net/http"
	"waldi-v2/controllers"
)

func main() {
	mux := http.NewServeMux()
	serveStaticFiles(mux)

	controllers.Init(mux)


	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func serveStaticFiles(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
