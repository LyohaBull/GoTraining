package main

import "fmt"

func newAccumulator(f float64) accumulator {
	return accumulator{f: f}
}

type accumulator struct {
	f float64
}

func (c *accumulator) read() {
	var b float64
	fmt.Scanf("%f", &b)
	c.f += b
}
func (c *accumulator) value() float64 {
	return c.f
}
func main() {
	c := newAccumulator(1.1)
	c.read()
	c.read()
	fmt.Println(c.value())
	c.read()
	fmt.Println(c.value())
}
