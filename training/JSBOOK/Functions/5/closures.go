package main

import (
	"fmt"
	"reflect"
	"time"
)

func sum(a ...int) {
	sum := 0
	for _, aa := range a {
		sum += aa
	}
	fmt.Println(sum)
}
func sumstr(a ...string) {
	fmt.Println(a)
}

var calls [][]reflect.Value

func spy(f any) func(a ...any) {
	ref := reflect.ValueOf(f)
	return func(a ...any) {
		inp := []reflect.Value{}
		for _, v := range a {
			inp = append(inp, reflect.ValueOf(v))
		}
		calls = append(calls, inp)
		ref.Call(inp)
	}
}

func deferred(f any, t time.Duration) func(a ...any) {
	ref := reflect.ValueOf(f)
	return func(a ...any) {
		inp := []reflect.Value{}
		for _, v := range a {
			inp = append(inp, reflect.ValueOf(v))
		}
		time.Sleep(t)
		ref.Call(inp)
	}
}
func kek() {
	fmt.Println("kek")
}

func main() {
	work := spy(sum)
	work(1, 2)
	work(3, 4, 5)
	work(5)
	work1 := spy(sumstr)
	work1("df", "dfdf")
	for _, a := range calls {
		fmt.Print("calls: ")
		for _, v := range a {
			fmt.Print(v, " ")
		}
		fmt.Println("")
	}

	delayedSum := deferred(sum, time.Second*3)
	delayedSum(1, 3)
	delayedSum(2, 5)
	delayedkek := deferred(kek, time.Second*10)
	delayedkek()
}
