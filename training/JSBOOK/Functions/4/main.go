package main

import (
	"fmt"
	"time"
)

func printNumbers(from, to int) {
	for i := from; i <= to; i++ {
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}
func main() {
	printNumbers(1, 5)
}
