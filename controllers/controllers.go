package controllers

import (
	"fmt"
	"net/http"
	h "waldi-v2/handlers"
	m "waldi-v2/handlers/middleware"
	s "waldi-v2/services"
)

func Init(mux *http.ServeMux) {
	appleMapsService := s.NewAppleMapsService()
	appleMapsClient := s.NewAppleMapsClient(appleMapsService)
	addressService := s.NewAddressService(appleMapsClient)

	baseRoute := New()
	baseRoute.SetParentRoute("/")
	baseRoute.RegisterRoute(Route{Method: "GET", Path: "", Handler: m.SetCacheHeader(h.Homepage)}, mux)

	resultRoute := New()
	resultRoute.SetParentRoute("/result")
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "", Handler: h.Resultpage(addressService)}, mux)
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "/", Handler: h.Resultpage(addressService)}, mux)
	resultRoute.RegisterRoute(Route{Method: "GET", Path: "/{id}", Handler: m.SetCacheHeader(h.Resultpage(addressService))}, mux)

	apiFormRoute := New()
	apiFormRoute.SetParentRoute("/api/form")
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/address", Handler: h.AddressForm(addressService)}, mux)
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/browserlocation", Handler: h.BrowserLocation(addressService)}, mux)
	apiFormRoute.RegisterRoute(Route{Method: "POST", Path: "/mapmodal", Handler: m.SetCacheHeader(h.MapModal())}, mux)
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
