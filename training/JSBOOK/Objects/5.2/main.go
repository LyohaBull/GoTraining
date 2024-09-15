package main

import "fmt"

func newCalculator() calculator {
	return calculator{}
}

type calculator struct {
	a float64
	b float64
}

func (c *calculator) read() {
	fmt.Scanf("%f %f", &c.a, &c.b)
}
func (c *calculator) sum() float64 {
	return c.a + c.b
}
func (c *calculator) mul() float64 {
	return c.a * c.b
}
func main() {
	c := newCalculator()
	c.read()
	fmt.Println(c.sum(), "  ", c.mul())
}
