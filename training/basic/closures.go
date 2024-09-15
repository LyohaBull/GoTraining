package main

import (
	"fmt"
	"strconv"
	"strings"
)

type my_int int
type hash_key string

func complexCalculating(arr ...my_int) my_int {
	sum := my_int(0)
	for i := range arr {
		sum += arr[i]
	}
	return sum
}

func hash(arr ...my_int) hash_key {
	sum := ""
	for i := range arr {
		sum += strconv.Itoa(int(arr[i])) + " "
	}
	sum = strings.TrimSpace(string(sum))
	return hash_key(sum)
}
func decorator(func(arr ...my_int) my_int) func(arr ...my_int) my_int {
	cache := make(map[hash_key]my_int, 10)
	return func(arr ...my_int) my_int {
		args := hash(arr...)
		if arg, ok := cache[args]; ok {
			fmt.Print("Закешировано: ")
			return arg
		} else {
			res := complexCalculating(arr...)
			cache[args] = res
			return res
		}

	}
}
func main() {
	complex := complexCalculating
	decorated := decorator(complex)

	i := 1
	for i != 0 {
		args := make([]my_int, 10)
		for j := 1; j != 0; {
			k := 0
			fmt.Scanf("%d ", &k)
			args = append(args, my_int(k))
			j = k
		}
		fmt.Println(decorated(args...))
	}
}
