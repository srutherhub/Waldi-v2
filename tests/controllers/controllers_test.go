package tests

import (
	"net/http"
	"testing"
	"waldi-v2/controllers"
)

func TestCreateController(t *testing.T) {
	testMux := http.NewServeMux()
	testHandler := func(w http.ResponseWriter, r *http.Request) {}
	testRoute := controllers.Route{Method: "GET", Path: "/test2", Handler: testHandler}

	testController := controllers.New()
	testController.SetParentRoute("/test")

	if testController.Base != "/test" {
		t.Fatal("SetParentRoute did not set the controllers base route")
	}

	testController.RegisterRoute(testRoute, testMux)

	if testController.Routes[0].Path != "/test2" {
		t.Fatal("RegisterRoute failed to register route")
	}
}
