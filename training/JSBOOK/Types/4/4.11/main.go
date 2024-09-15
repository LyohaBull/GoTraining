package main

import (
	"fmt"
	"slices"
)

func main() {
	strs := []string{"кришна", "кришна", "харе", "харе",
		"харе", "харе", "кришна", "кришна", ":-O"}
	slices.Compact(strs)
	fmt.Println(strs)
}
