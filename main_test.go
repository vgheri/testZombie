//main_test.go
package main

import "testing"

//TestAlgo ensure that the formula used returns the correct result
func TestAlgo(t *testing.T) {
	lat1, lon1 := 48.856614, 2.352222 //Paris
	lat2, lon2 := 43.710173, 7.261953 //Nice

	r := Distance(lat1, lon1, lat2, lon2)

	if r != 685872.0703773521 {
		t.Fail()
	}
}

//TestisDriverZombie check that our methods does well the sum of distance
func TestIsZombieShouldReturnNotZombieForParisNiceParis(t *testing.T) {
	var locations = []*DriverLocation{
		&DriverLocation{
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: "",
		},
		&DriverLocation{
			Latitude:  43.710173,
			Longitude: 7.261953,
			UpdatedAt: "",
		},
		&DriverLocation{
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: "",
		},
	}
	r := isDriverZombie(locations)

	if r == true {
		t.Fail()
	}
}

func TestIsZombieShouldReturnIsZombieForSingleLocation(t *testing.T) {
	var locations = []*DriverLocation{
		&DriverLocation{
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: "",
		},
	}
	r := isDriverZombie(locations)

	if r == false {
		t.Fail()
	}
}

func TestIsZombieShouldReturnIsZombieForTwoLocationsApartLessThan500m(t *testing.T) {
	var locations = []*DriverLocation{
		&DriverLocation{
			Latitude:  48.8566,
			Longitude: 2.3522,
			UpdatedAt: "",
		},
		&DriverLocation{
			Latitude:  48.8566,
			Longitude: 2.3523,
			UpdatedAt: "",
		},
	}
	r := isDriverZombie(locations)

	if r == false {
		t.Fail()
	}
}
