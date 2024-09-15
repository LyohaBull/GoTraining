package main

import "fmt"

type person struct {
	name string
	age  int
}

func main() {
	vasya := person{"vasya", 25}
	petya := person{"petya", 28}
	masha := person{"masha", 22}
	arr := []person{vasya, petya, masha}
	names := make([]string, len(arr))
	for i := range arr {
		names[i] = arr[i].name
	}
	fmt.Println(names)
}
