package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

func aclean(arr []string) []string {
	m := make(map[string]bool, len(arr))
	del := make([]bool, len(arr))
	for i, str := range arr {
		str = strings.ToLower(str)
		r := make([]rune, len(str))
		for j, s := range str {
			r[j] = s

		}
		slices.Compact(r)
		slices.Sort(r)
		ss := ""
		for _, s := range r {
			ss = string(utf8.AppendRune([]byte(ss), s))
		}
		if m[ss] {
			del[i] = false
		} else {
			del[i] = true
			m[ss] = true
		}
	}
	slices.DeleteFunc(arr, func(n string) bool {
		return !del[slices.Index(arr, n)]
	})
	return arr[:len(m)]
}

func main() {
	arr := []string{"Hare", "Krishna", "Hare", "Krishna",
		"Krishna", "Krishna", "Hare", "Hare", ":-O"}
	set := make(map[string]bool)
	for _, str := range arr {
		set[str] = true
	}
	fmt.Println(set)

	strs := []string{"nap", "teachers", "cheaters", "PAN", "ear", "era", "hectares"}
	strs = aclean(strs)
	fmt.Print(strs)

}
