package main

import "fmt"

type ladder int

func (l *ladder) showStep() *ladder {
	fmt.Println(*l)
	return l
}
func (l *ladder) up() *ladder {
	*l++
	return l
}
func (l *ladder) down() *ladder {
	*l--
	return l
}
func main() {
	ladder := ladder(0)
	ladder.up().up().down().showStep().down().showStep() // показывает 1 затем 0
}
