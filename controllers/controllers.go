package controllers

import (
	"fmt"
	"net/http"
	"waldi-v2/handlers"
)

func Init(mux *http.ServeMux) {
	baseRoute := New()
	baseRoute.setParentRoute("/")
	baseRoute.registerRoute(Route{Method: "GET", Path: "", Handler: handlers.Homepage}, mux)

	apiRoute := New()
	apiRoute.setParentRoute("/api")
	apiRoute.registerRoute(Route{Method: "POST", Path: "/submitlocation", Handler: handlers.SubmitLocation},mux)
}

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Controller struct {
	base   string
	routes []Route
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) setParentRoute(path string) {
	c.base = path
}

func (c *Controller) registerRoute(route Route, mux *http.ServeMux) {
	mux.HandleFunc(c.base+route.Path, route.Handler)
	fmt.Println("Registered: " + c.base+route.Path)
}
