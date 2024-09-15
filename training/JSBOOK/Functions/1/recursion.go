package main

import (
	"fmt"
	"slices"
	"strconv"
)

func sumTo1(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}
func sumTo2(n int) int {
	if n == 2 {
		return 3
	}
	return n + sumTo2(n-1)
}
func sumTo3(n int) int {
	return ((1 + n) * n / 2)
}
func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * factorial(n-1)
}
func fib(n int) int {
	a := 1
	b := 1
	for i := 3; i <= n; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}

type list struct {
	value int
	next  *list
}

func printList(l list) string {
	if l.next == nil {
		return strconv.Itoa(l.value)
	}
	return strconv.Itoa(l.value) + " " + printList(*l.next)
}
func printList1(l list) {
	for l.next != nil {
		fmt.Print(l.value, " ")
		l = *l.next
	}
	fmt.Print(l.value, "\n")
}
func printOutList1(l list) {
	arr := []int{}
	for l.next != nil {
		arr = append(arr, l.value)
		l = *l.next
	}
	arr = append(arr, l.value)
	slices.Reverse(arr)
	for _, a := range arr {
		fmt.Print(a, " ")
	}
	fmt.Println("")
}
func printOutList(l list) string {
	if l.next == nil {
		return strconv.Itoa(l.value)
	}
	return printOutList(*l.next) + " " + strconv.Itoa(l.value)
}
func main() {
	fmt.Println(sumTo1(100), " ", sumTo2(100), " ", sumTo3(100))
	fmt.Println(factorial(5))
	fmt.Println(fib(77))
	list4 := list{4, nil}
	list3 := list{3, &list4}
	list2 := list{2, &list3}
	list1 := list{1, &list2}
	fmt.Println(printList(list1))
	printList1(list1)
	printOutList1(list1)
	fmt.Println(printOutList(list1))

}
