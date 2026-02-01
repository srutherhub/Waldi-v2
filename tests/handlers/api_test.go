package tests

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

var Url = "http://localhost:8080"

func TestSubmitLocation(t *testing.T) {
	address := "Madison Wisconsin"

	formData := url.Values{}
	formData.Set("address", address)

	resp, err := http.PostForm(Url+"/api/form/address", formData)

	if err != nil {
		t.Fatal("/api/form/address api call failed: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("/api/form/address failed to read body: " + err.Error())
	}

	expected := `<form hx-post="api/form/address" hx-target="this" hx-swap="outerHTML" class="hstack gap1"><div class="w100 relative inlineblock"><input type="text" name="address" placeholder="Address" class="input" value="Madison Wisconsin"><div class="absolute front" style="right: 5px; top: 50%; transform: translateY(-50%);"><button class="button2" type="">My Location</button></div></div><button class="button" type="submit">Submit</button></form>`

	// fmt.Println(string(body))
	// fmt.Println(expected)

	if string(body) != expected {
		t.Fatal("/api/form/address response does not match expected html")
	}

}
