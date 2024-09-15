package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ucFirst(s string) string {
	ss, _ := utf8.DecodeRuneInString(s)
	return strings.Replace(s, string(ss), strings.ToUpper(string(ss)), 1)
}
func truncate(s string, size int) string {
	col := utf8.RuneCountInString(s)
	if col <= size {
		return s
	}
	i := 19
	for {
		col1 := utf8.RuneCountInString(s[:i])
		if col1 == 20 {
			break
		}
		i++
	}
	i--
	return s[:i] + string('\u2026')
}
func extractCurrencyValue(s string) float64 {
	_, size := utf8.DecodeRuneInString(s)
	f, err := strconv.ParseFloat(s[size:], 64)
	if err != nil {
		fmt.Println("err")
		return 0
	}
	return f
}

func main() {
	str := "лёха"
	fmt.Println(ucFirst(str))
	fmt.Println(strings.Contains(strings.ToLower("Привет лёха!"), strings.ToLower("Лёха")))
	fmt.Println(truncate("Вот, что мне хотелось бы сказать на эту тему:", 20))
	fmt.Println(truncate("Всем привет!", 20))
	fmt.Println(extractCurrencyValue("$12.1"))
}
