package main

import (
	"fmt"
	"slices"
)

type person struct {
	id   string
	name string
	age  int
}

func sortByAge(arr []person) {
	slices.SortFunc(arr, func(a, b person) int {
		return a.age - b.age
	})
}
func getAverageAge(arr []person) float64 {
	average := 0.0
	for _, person := range arr {
		average += float64(person.age)
	}
	return average / float64(len(arr))
}

type persons map[string]person

func groupBy(arr []person) persons {
	m := persons{}
	for _, person := range arr {
		m[person.id] = person
	}
	return m
}
func main() {
	users := []person{
		{id: "john", name: "John Smith", age: 20},
		{id: "ann", name: "Ann Smith", age: 24},
		{id: "pete", name: "Pete Peterson", age: 31},
	}
	sortByAge(users)
	fmt.Println(users)
	fmt.Println(getAverageAge(users))
	fmt.Println(groupBy(users))
}
