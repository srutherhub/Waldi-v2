package tests

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

var Url = "http://localhost:8080"

func TestAddressForm(t *testing.T) {
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

	expected := `<form hx-post="api/form/address" hx-target="this" hx-swap="outerHTML" id="form-address" class="hstack gap1"><div class="w100 relative inlineblock"><input type="text" name="address" id="inp.address" placeholder="Address" class="input" value="Madison Wisconsin"><div class="absolute front" style="right: 5px; top: 50%; transform: translateY(-50%);"><button class="button2" id="btn-mylocation" type="button">My Location</button></div></div><button class="button" id="btn-submitaddress">Submit</button></form><script>
  addButtonOnClick("btn-mylocation", myLocationHandler)
</script>`

	// fmt.Println(string(body))
	// fmt.Println(expected)

	if string(body) != expected {
		t.Fatal("/api/form/address response does not match expected html")
	}

}

func TestBrowserLocation(t *testing.T) {
	lat := "43"
	lon := "-89"

	formData := url.Values{}
	formData.Set("lat", lat)
	formData.Set("lon", lon)

	resp, err := http.PostForm(Url+"/api/form/browserlocation", formData)

	if err != nil {
		t.Fatal("/api/form/address api call failed: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("/api/form/address failed to read body: " + err.Error())
	}

	expected := `<form hx-post="api/form/address" hx-target="this" hx-swap="outerHTML" id="form-address" class="hstack gap1"><div class="w100 relative inlineblock"><input type="text" name="address" id="inp.address" placeholder="Address" class="input" value="43, -89"><div class="absolute front" style="right: 5px; top: 50%; transform: translateY(-50%);"><button class="button2" id="btn-mylocation" type="button">My Location</button></div></div><button class="button" id="btn-submitaddress">Submit</button></form><script>
  addButtonOnClick("btn-mylocation", myLocationHandler)
</script>`

	// fmt.Println(string(body))
	// fmt.Println(expected)

	if string(body) != expected {
		t.Fatal("/api/form/address response does not match expected html")
	}

}