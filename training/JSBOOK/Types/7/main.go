package main

import (
	"cmp"
	"fmt"
	"slices"
)

type person struct {
	name   string
	salary int
}

func topSalary(arr []person) int {
	return slices.MaxFunc(arr, func(a, b person) int {
		return cmp.Compare(a.salary, b.salary)
	}).salary
}
func main() {
	arrs := []person{
		{"John", 100},
		{"Pete", 300},
		{"Mary", 250},
	}
	fmt.Println(topSalary(arrs))
}
