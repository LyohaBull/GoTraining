package main

import "fmt"

type salary struct {
	name   string
	salary int
}

func getSalaries(s []salary) int {
	sum := 0
	for _, salary := range s {
		sum += salary.salary
	}
	return sum
}
func main() {
	org := []salary{
		{"Alice", 100},
		{"Bob", 200},
		{"John", 300},
		{"Alice", 100},
		{"Bob", 200},
		{"John", 300},
		{"Alice", 100},
		{"Bob", 200},
		{"John", 300},
	}
	fmt.Println(getSalaries(org))
}
