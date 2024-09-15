package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type person struct {
	Name string `json:"my_name"`
	Age  int
}

func main() {
	john := person{"John", 30}
	arr := []person{john, {"Ann", 45}, {"Jack", 21}}
	bytes, err := json.Marshal(arr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	bytes2, err := json.Marshal(john)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes2))
	var new_john person
	err1 := json.Unmarshal(bytes2, &new_john)
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Printf("%+v\n", new_john)
}
