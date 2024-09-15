package main

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
)

func sum(a int) func(int) int {
	return func(b int) int {
		return a + b
	}
}

func inBetween[T cmp.Ordered](arr []T, a, b T) []T {
	newArr := []T{}
	for _, m := range arr {
		if m >= a && m <= b {
			newArr = append(newArr, m)
		}
	}
	return newArr
}

type person struct {
	name    string
	age     int
	surname string
}

func byField(s string) func(person, person) int {
	return func(a, b person) int {
		t1, _ := reflect.TypeOf(a).FieldByName(s)
		v1 := reflect.ValueOf(a)
		v2 := reflect.ValueOf(b)
		fmt.Println(t1.Name)
		switch t1.Type.Kind().String() {
		case "int":
			{
				return cmp.Compare(v1.FieldByName(s).Int(), v2.FieldByName(s).Int())
			}
		case "string":
			{
				return cmp.Compare(v1.FieldByName(s).String(), v2.FieldByName(s).String())
			}
		default:
			{
				fmt.Println("default")
				return cmp.Compare(v1.FieldByName(s).Int(), v2.FieldByName(s).Int())
			}
		}
	}
}

func main() {
	fmt.Println(sum(102)(104))
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(inBetween(arr, 3, 6))
	users := []person{
		{name: "John", age: 20, surname: "Johnson"},
		{name: "Pete", age: 18, surname: "Peterson"},
		{name: "Ann", age: 19, surname: "Hathaway"},
	}
	slices.SortFunc(users, byField("age"))
	fmt.Println(users)

}
