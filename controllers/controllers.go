package controllers

import (
	"fmt"
	"net/http"
	"waldi-v2/handlers"
	"waldi-v2/services"
)

func Init(mux *http.ServeMux) {
	appleMapsService := services.NewAppleMapsService()
	appleMapsClient := services.NewAppleMapsClient(appleMapsService)
	addressService := services.NewAddressService(appleMapsClient)

	baseRoute := New()
	baseRoute.SetParentRoute("/")
	baseRoute.RegisterRoute(Route{Method: "GET", Path: "", Handler: handlers.Homepage}, mux)

	resultRoute := New()
	resultRoute.SetParentRoute("/result")
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "", Handler: handlers.Resultpage(addressService)}, mux)
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "/", Handler: handlers.Resultpage(addressService)}, mux)
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "/{id}", Handler: handlers.Resultpage(addressService)}, mux)

	apiFormRoute := New()
	apiFormRoute.SetParentRoute("/api/form")
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/address", Handler: handlers.AddressForm(addressService)}, mux)
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/browserlocation", Handler: handlers.BrowserLocation(addressService)}, mux)
		apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/mapmodal", Handler: handlers.MapModal()}, mux)
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
	mux.HandleFunc(c.Base+route.Path, route.Handler)

	c.Routes = append(c.Routes, route)

	fmt.Println("Registered: " + c.Base + route.Path)
}
