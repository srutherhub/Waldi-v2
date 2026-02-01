package controllers

import (
	"fmt"
	"net/http"
	"waldi-v2/handlers"
)

func Init(mux *http.ServeMux) {
	baseRoute := New()
	baseRoute.SetParentRoute("/")
	baseRoute.RegisterRoute(Route{Method: "GET", Path: "", Handler: handlers.Homepage}, mux)

	apiFormRoute := New()
	apiFormRoute.SetParentRoute("/api/form")
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/address", Handler: handlers.AddressForm}, mux)
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/browserlocation", Handler: handlers.BrowserLocation}, mux)
}

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Controller struct {
	Base   string
	Routes []Route
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) SetParentRoute(path string) {
	c.Base = path
}

func (c *Controller) RegisterRoute(route Route, mux *http.ServeMux) {
	mux.HandleFunc(c.Base + route.Path, route.Handler)

	c.Routes = append(c.Routes, route)

	fmt.Println("Registered: " + c.Base + route.Path)
}
