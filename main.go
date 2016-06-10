package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/drivers/{id:[0-9]+}", ZombieDriverHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Server started and listening on port %d.", 1338)
	log.Fatal(http.ListenAndServe(":1338", nil))
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

	locations, err := getDriverLocations(driverID)
	if err != nil {
		// TODO add circuit breaker system
		log.Printf("Could not retrieve driver's locations.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isZombie := isDriverZombie(locations)

	result := APIResponse{ID: driverID, Zombie: isZombie}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(result)
	if err != nil {
		log.Printf("Could not json encode response.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	return
}

func getDriverLocations(id int) ([]*DriverLocation, error) {
	return []*DriverLocation{
		&DriverLocation{
			Latitude:  42,
			Longitude: 2.3,
			UpdatedAt: "2016-06-10T19:43:22.232Z",
		},
		&DriverLocation{
			Latitude:  42,
			Longitude: 2.3,
			UpdatedAt: "2016-06-10T19:43:22.232Z",
		},
	}, nil
}

func isDriverZombie(locations []*DriverLocation) bool {
	return false
}
