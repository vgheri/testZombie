//main_test.go
package main

import (
	"fmt"
	"testing"
)

//TestAlgo ensure that the formula used returns the correct result
func TestAlgo(t *testing.T) {
	lat1, lon1 := 48.856614, 2.352222 //Paris
	lat2, lon2 := 43.710173, 7.261953 //Nice

	r := Distance(lat1, lon1, lat2, lon2)
	fmt.Println(r)

	if r != 685872.0703773521 {
		t.Fail()
	}
}
