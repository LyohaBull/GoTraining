package main

import (
	"fmt"
	"strconv"
	"time"
)

func test(c chan int) {

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(int(time.Second)))
		fmt.Print(strconv.Itoa(i+1) + ": ")
		c <- i + 1
	}
	close(c)
}
func plus(c, c2 chan int) {
	for {
		a, ok := <-c
		if !ok {
			close(c2)
			return
		}
		c2 <- a + 10*a
	}

}
func result(c chan int) {
	for {
		res, ok := <-c
		if !ok {
			return
		}
		fmt.Println(res)
	}
}

func main() {
	c := make(chan int)
	c2 := make(chan int)
	go test(c)
	go plus(c, c2)
	result(c2)
}
