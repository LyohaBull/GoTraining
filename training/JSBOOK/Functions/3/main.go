package main

import (
	"fmt"
	"strconv"
)

type counter struct {
	n int
}

func (c *counter) counter() int {
	c.n++
	return c.n
}
func (c *counter) set(a int) int {
	c.n = a
	return c.n
}
func (c *counter) decrease() int {
	c.n--
	return c.n
}
func makeCounter() counter {
	return counter{1}
}

type su func(int) su

func (s su) String() string {
	return strconv.Itoa(currentSum)
}

var currentSum int

func f(b int) su {
	currentSum += b
	return f
}

func sum(a int) su {
	currentSum = a
	s := su(f)
	return s
}

func main() {
	c := makeCounter()
	fmt.Println(c.counter())
	fmt.Println(c.set(20))
	fmt.Println(c.decrease())
	fmt.Println(c.counter())
	fmt.Println(sum(1)(2)(3)(20))
	fmt.Println(sum(1)(2)(3)(20))
}
