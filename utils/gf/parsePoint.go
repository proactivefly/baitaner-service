package gf

import (
	"fmt"
)

type Point struct {
	Latitude  float64
	Longitude float64
}

func parsePoint(wkt string) (*Point, error) {
	fmt.Println("parsePoint")
	return nil, nil
}
