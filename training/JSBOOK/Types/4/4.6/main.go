package main

import (
	"fmt"
	"strconv"
	"strings"
)

type calculator struct {
	f map[string]func(a, b int) int
}

func (c calculator) calculate(str string) int {
	arr := strings.Split(str, " ")
	a, _ := strconv.Atoi(arr[0])
	b, _ := strconv.Atoi(arr[2])
	return c.f[arr[1]](a, b)
}

func (c calculator) addMethod(str string, method func(a int, b int) int) {
	c.f[str] = method
}

func main() {
	c := calculator{f: map[string]func(a, b int) int{
		"+": func(a, b int) int {
			return a + b
		},
		"-": func(a, b int) int {
			return a - b
		},
	}}
	fmt.Println(c.calculate("1 - 2"))
	c.addMethod("*", func(a, b int) int {
		return a * b
	})
	c.addMethod("/", func(a, b int) int {
		return a / b
	})
	fmt.Println(c.calculate("4 / 2"))
}
