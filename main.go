package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	consul "github.com/hashicorp/consul/api"
	"github.com/rubyist/circuitbreaker"
)

var breaker = circuit.NewConsecutiveBreaker(5)
var consulClient *consul.Client

// DriverLocation models response from LocationService
type DriverLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

// APIResponse models the response to know if a driver is active
type APIResponse struct {
	ID     int  `json:"id"`
	Zombie bool `json:"zombie"`
}

func main() {
	var err error
	if consulClient, err = initConsul(); err != nil {
		log.Fatal(err)
	}
	err = register(consulClient, "zombie", "172.17.0.1", 1338)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/drivers/{id:[0-9]+}", ZombieDriverHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Server started and listening on port %d.", 1338)
	log.Println(http.ListenAndServe(":1338", nil))
	unregister(consulClient, "zombie")
}

// ZombieDriverHandler handles a user's request to know if a driver is active
func ZombieDriverHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("\t%s",
		r.RequestURI)
	// Read route parameter
	vars := mux.Vars(r)
	param := vars["id"]
	driverID, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Received bad request with driver id %s.", param)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	baseAddr, err := retrieveLocationServiceAddress(consulClient)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	locationURL := fmt.Sprintf("http://%s/drivers/%d/coordinates?minutes=5", baseAddr, driverID)
	locations, err := getDriverLocations(breaker, locationURL)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	isZombie := isDriverZombie(locations)
	result := APIResponse{ID: driverID, Zombie: isZombie}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(result)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	return
}

func getDriverLocations(breaker *circuit.Breaker, serviceURL string) ([]*DriverLocation, error) {
	var response *http.Response
	var err error
	err = breaker.Call(func() error {
		var httpErr error
		response, httpErr = http.Get(serviceURL)
		return httpErr
	}, time.Second*1)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	var locations []*DriverLocation
	if err := decoder.Decode(&locations); err != nil {
		return nil, err
	}
	return locations, nil
}

//isDriverZombie
func isDriverZombie(locations []*DriverLocation) bool {
	var meters float64
	var prev DriverLocation
	for i, loc := range locations {
		if i > 0 {
			meters = Distance(prev.Latitude, prev.Longitude, loc.Latitude, loc.Longitude) + meters
		}
		prev.Latitude = loc.Latitude
		prev.Longitude = loc.Longitude
	}

	if meters < 500 {
		return true
	}
	return false
}

func retrieveLocationServiceAddress(client *consul.Client) (string, error) {
	return service(consulClient, "location", "")
}
