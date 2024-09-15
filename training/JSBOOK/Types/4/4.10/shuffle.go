package main

import (
	"fmt"
	"math/rand"
)

func main() {
	arr := []int{1, 2, 3}
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	fmt.Println(arr)
}
