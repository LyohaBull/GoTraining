package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode/utf8"
)

func sumInput() float64 {
	sum := 0.0
	for {

		f := 0.0
		_, err := fmt.Scanf("%f", &f)
		if err != nil {
			fmt.Println("error in scan")
			return sum
		}
		sum += f
	}
}
func getMaxSubArr(s []float64) float64 {
	maxSum := 0.0
	partialSum := 0.0
	for _, a := range s {
		partialSum += a
		maxSum = math.Max(maxSum, partialSum)
		if partialSum < 0 {
			partialSum = 0
		}
	}
	return maxSum
}
func ucFirst(s string) string {
	ss, _ := utf8.DecodeRuneInString(s)
	return strings.Replace(s, string(ss), strings.ToUpper(string(ss)), 1)
}
func camelize(s string) string {
	a := strings.Split(s, "-")
	res := ""
	for i := range a {
		if i == 0 {
			res += a[i]
			continue
		}
		res += ucFirst(a[i])
	}
	return res
}
func filterRange(arr []int, a, b int) []int {
	res := []int{}
	for _, r := range arr {
		if a <= r && b >= r {
			res = append(res, r)
		}
	}
	return res
}
func filterRangeInPlace(arr []int, a, b int) {
	slices.DeleteFunc(arr, func(n int) bool {
		return n < a || n > b
	})
}
func copySorted(arr []string) []string {
	var res []string
	res = append(res, arr...)
	slices.Sort(res)
	return res
}
func main() {
	//fmt.Println(sumInput())
	fmt.Println(getMaxSubArr([]float64{100, -9, 2, -3, 5}))
	fmt.Println(camelize("-webkit-transition"))
	fmt.Println(filterRange([]int{5, 3, 8, 1}, 1, 4))
	aa := []int{5, 3, 8, 1}
	filterRangeInPlace(aa, 1, 4)
	fmt.Println(aa)
	bb := []int{5, 2, 1, -10, 8}
	slices.SortFunc(bb, func(a, b int) int {
		return b - a
	})
	fmt.Println(bb)
	cc := []string{"HTML", "JavaScript", "CSS"}
	fmt.Println(copySorted(cc))
	fmt.Println(cc)
}
