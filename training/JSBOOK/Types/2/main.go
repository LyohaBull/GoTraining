package main

import (
	"fmt"
	"math/rand"
)

func randMinMax(min, max float64) float64 {
	return min + (max-min)*rand.Float64()

}
func main() {
	a := randMinMax(-23.2, 1102.2123)
	fmt.Println(a)

}
