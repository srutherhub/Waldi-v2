# App Architecture

## Components/Views
Components (or views) determine the UI. Made with go templ.

## Controllers
Controllers determine what endpoints are and what each endpoint does.

## Handlers
Handlers takes http requests and responds to them with data.

## Services
Services hold core app logic and are passed into handlers or other services. They are independent of the web layer (handlers/controllers).
ex. AppleMapsService is passed into AppleMapsClient. AppleMapsClient is passed into AddressService. AddressService is passed into a handler.

## Clients
Are similar to services. Clients interact with outside sources like APIs, and are defined in a way where they can be interchanged.
ex. AddressService requires an IAddressClient to function. IAddressClient can do two things, geocode and reverse geocode. As long as 
any client (Apple Maps, Google, MapBox) etc. implements geocoding and reverse geocoding, it can be plugged into AddressService.