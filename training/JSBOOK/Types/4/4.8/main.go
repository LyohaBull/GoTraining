package main

import "fmt"

type person struct {
	name    string
	surname string
	id      int
}
type user struct {
	fullname string
	id       int
}

func (p person) toUser() user {
	return user{fullname: p.name + " " + p.surname, id: p.id}
}
func main() {
	vasya := person{"Вася", "Пупкин", 1}
	fmt.Println(vasya.toUser())
}
