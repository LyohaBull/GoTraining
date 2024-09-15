package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type action struct {
	Pref string
	Num  int
	GrPr string
}
type macrosline struct {
	n    int
	nn   int
	line out_map
}

func (a action) String() string {
	return a.Pref + "," + fmt.Sprintf("%03d", a.Num) + " - 0 - " + a.GrPr
}

type Num struct {
	input     int
	Num_start string
	Num_end   string
	act       action
}
type Num_c struct {
	input int
	Num   string
	act   action
}

func (n Num) toNumC() []Num_c {
	Nums := []Num_c{}
	nst, err := strconv.Atoi(n.Num_start)
	if err != nil {
		log.Fatal("err in Num_start")
		return []Num_c{}
	}
	ned, err := strconv.Atoi(n.Num_end)
	if err != nil {
		log.Fatal("err in Num_start")
		return []Num_c{}
	}
	diff := ned - nst + 1
	i := len(n.Num_start) - 1
	sm := 0
	rev := false
	for diff > 0 {
		f := 0
		if diff < int(math.Pow10(sm)) && !rev {
			rev = true
			i += 1
			sm -= 1
			continue
		}
		if rev {
			if diff < int(math.Pow10(sm)) {
				i++
				sm--
				continue
			}
			f = diff/int(math.Pow10(sm)) - 1
			ned, err = strconv.Atoi(string(n.Num_end[i]))
			if err != nil {
				log.Fatal("err")
				return []Num_c{}
			}
			end := int(math.Min(float64(0+f), float64(ned)))
			for l := 0; l <= end; l++ {
				Nums = append(Nums, Num_c{n.input, string(n.Num_start[:i]) + strconv.Itoa(l), n.act})
				n.Num_start = string(n.Num_start[:i]) + strconv.Itoa(l) + string(n.Num_start[i+1:])
				diff -= int(math.Pow10(sm))
			}
			n.Num_start = string(n.Num_start[:i]) + strconv.Itoa(end+1) + string(n.Num_start[i+1:])
			if diff <= 0 {
				break
			}
			i++
			sm--
		} else {
			if string(n.Num_start[i]) == "0" && diff > 9 && i != 0 {
				i--
				sm++
				continue
			}
			f = diff/int(math.Pow10(sm)) - 1
			nst, err = strconv.Atoi(string(n.Num_start[i]))
			if err != nil {
				log.Fatal("err")
				return []Num_c{}
			}
			end := int(math.Min(float64(nst+f), float64(int((10*(sm+1)-1)%10))))

			for l := nst; l <= end; l++ {
				Nums = append(Nums, Num_c{n.input, string(n.Num_start[:i]) + strconv.Itoa(l), n.act})
				n.Num_start = string(n.Num_start[:i]) + strconv.Itoa(l) + string(n.Num_start[i+1:])
				diff -= int(math.Pow10(sm))
			}
			if diff <= 0 {
				break
			}
			if end == 9 {
				n.Num_start = string(n.Num_start[:i]) + "0" + string(n.Num_start[i+1:])
				h := 9
				var err error
				for j := 1; h == 9; j++ {
					h, err = strconv.Atoi(string(n.Num_start[i-j]))
					if err != nil {
						log.Fatal("err")
						return []Num_c{}
					}
					n.Num_start = string(n.Num_start[:i-j]) + strconv.Itoa((h+1)%10) + string(n.Num_start[i-j+1:])
				}
				i--
				sm++
			} else {
				n.Num_start = string(n.Num_start[:i]) + strconv.Itoa(end+1) + string(n.Num_start[i+1:])
			}
		}
	}
	return Nums

}

func (n Num_c) getMap(index int, Nums *[][]out_map) (int, error) {
	k := n.input
	for i := 0; i < len(n.Num); i++ {
		nn := n.getNum(i)
		line := (*Nums)[k][nn]
		if line.Transition {
			if i == index {
				return k, nil
			}
			k = line.Transit
			continue
		}
		if line.Act_ptr != nil {
			return k, nil
		}
		return 0, errors.New("нет такого вызова")
	}
	return 0, errors.New("нет такого вызова")
}

type out_map struct {
	Transit    int
	Transition bool
	Act_ptr    *action
}
type map_table struct {
	Nums *[][]out_map
	Next *[]int
}

func (n Num_c) getCall(Nums *[][]out_map) (action, []Num_c, error) {
	k := n.input
	for i := 0; i < len(n.Num); i++ {
		nn := n.getNum(i)
		line := (*Nums)[k][nn]
		if line.Transition {
			k = line.Transit
			continue
		}
		if line.Act_ptr != nil {
			return *line.Act_ptr, []Num_c{}, nil
		}
		return action{}, []Num_c{}, errors.New("нет такого вызова")
	}
	new := make([]Num_c, 10)
	for i := 0; i < 10; i++ {
		new[i] = Num_c{n.input, n.Num + strconv.Itoa(i), n.act}
	}
	return action{}, new, nil
}

func (mt *map_table) upgrade(n Num_c) ([]int, []macrosline) {
	map_tbls := []int{}
	lines := []macrosline{}
	if len(n.Num) == 1 {
		return map_tbls, lines
	}
	lastmap, err := n.getMap(len(n.Num)-1, mt.Nums)
	if err != nil {
		log.Fatal("err in getMap")
		return map_tbls, lines
	}
	l := (*mt.Nums)[lastmap][0].Act_ptr
	if l == nil {
		return map_tbls, lines
	}
	for _, line := range (*mt.Nums)[lastmap] {
		if line.Act_ptr == nil {
			return map_tbls, lines
		}
		if *line.Act_ptr != *l {
			return map_tbls, lines
		}
	}
	var prevlastmap int
	prevlastmap, err = n.getMap(len(n.Num)-2, mt.Nums)
	if err != nil {
		log.Fatal("err in getMap")
		return map_tbls, lines
	}
	for i, line := range (*mt.Nums)[prevlastmap] {
		if line.Transit == lastmap && line.Transition {
			(*mt.Nums)[prevlastmap][i].Transition = false
			(*mt.Nums)[prevlastmap][i].Act_ptr = &n.act
			lines = append(lines, macrosline{prevlastmap, i, line})
			break
		}
		if i == 9 {
			log.Fatal("err in upgrade")
			return map_tbls, lines
		}
	}

	(*mt.Nums)[lastmap] = make([]out_map, 10)
	for i, line := range (*mt.Nums)[lastmap] {
		lines = append(lines, macrosline{lastmap, i, line})
	}
	map_tbls = append(map_tbls, prevlastmap)
	map_tblsU, linesU := mt.upgrade(Num_c{n.input, n.Num[:len(n.Num)-1], n.act})
	map_tbls = append(map_tbls, map_tblsU...)
	lines = append(lines, linesU...)
	if (len(*mt.Next)) == 0 {
		(*mt.Next) = append((*mt.Next), lastmap)

		return map_tbls, lines
	} else {
		(*mt.Next) = append((*mt.Next), lastmap)
		slices.Sort((*mt.Next))
	}
	return map_tbls, lines
}

func getCalls(n []Num_c, Nums *[][]out_map) ([]Num_c, []action) {
	res := []action{}
	res_nums := []Num_c{}
	for _, Num := range n {
		a, new, err := Num.getCall(Nums)
		if err != nil {
			res = append(res, a)
			res_nums = append(res_nums, Num)
			continue
		}
		if len(new) == 0 {
			res = append(res, a)
			res_nums = append(res_nums, Num)
		} else {
			add_nums, add_res := getCalls(new, Nums)
			res = append(res, add_res...)
			res_nums = append(res_nums, add_nums...)
		}
	}
	return res_nums, res
}
func printCalls(filename string, n []Num_c, a []action) {
	var f *os.File
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		for i, act := range a {
			if act.String() != ",000 - 0 - " {
				io.WriteString(f, strconv.Itoa(n[i].input)+":"+n[i].Num+" - "+act.String()+"\r\n")
				continue
			}
			io.WriteString(f, strconv.Itoa(n[i].input)+":"+n[i].Num+" - нет такого вызова"+"\r\n")
		}
		fmt.Println("Файл " + filename + " успешно записан!")
	} else {
		fmt.Println("Файл с именем " + filename + " уже существует. Заменить его? y/n")
		yn := ""
		fmt.Scanf("%s\r\n", &yn)
		if yn == "y" {
			os.Remove(filename)
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			for i, act := range a {
				if act.String() != ",000 - 0 - " {
					io.WriteString(f, strconv.Itoa(n[i].input)+":"+n[i].Num+" - "+act.String()+"\r\n")
					continue
				}
				io.WriteString(f, strconv.Itoa(n[i].input)+":"+n[i].Num+" - нет такого вызова"+"\r\n")
			}
			fmt.Println("Файл " + filename + " успешно записан!")
		} else {
			fmt.Println("Новый файл " + filename + " не записан!")
		}
	}

}

func readline(s string, f bool) (Num, error) {
	arr := strings.Split(s, ":")
	if len(arr) < 3 && f {
		return Num{}, errors.New("неправильные входные данные")
	}
	var n Num
	var err error
	n.input, err = strconv.Atoi(arr[0])
	if err != nil {
		return Num{}, errors.New("недопустимый вход в ПН")
	}
	if strings.Contains(arr[1], "-") {
		arr1 := strings.Split(arr[1], "-")
		if cmp.Compare(arr1[0], arr1[1]) == 1 {
			return Num{}, errors.New("конец диапазона < начала диапазона")
		}
		_, err := strconv.Atoi(arr1[0])
		if err != nil {
			return Num{}, errors.New("недопустимое число начала диапазона")
		}
		_, err = strconv.Atoi(arr1[1])
		if err != nil {
			return Num{}, errors.New("недопустимое число конца диапазона")
		}
		n.Num_start = arr1[0]
		n.Num_end = arr1[1]
	} else {
		_, err := strconv.Atoi(arr[1])
		if err != nil {
			return Num{}, errors.New("недопустимое число номера")
		}
		n.Num_start = arr[1]
		n.Num_end = n.Num_start
	}
	if f {
		if strings.Contains(arr[2], ",") {
			arr1 := strings.Split(arr[2], ",")
			act := action{arr1[0], 0, "-"}
			act.Num, err = strconv.Atoi(arr1[1])
			if err != nil {
				return Num{}, errors.New("ошибка в вызове")
			}
			if len(arr1) > 2 {
				if arr1[2] == "+" {
					act.GrPr = arr1[2]
				}
			}
			n.act = act
		} else {
			return Num{}, errors.New("ошибка в вызове")
		}
	} else {
		n.act = action{}
		return n, nil
	}
	return n, nil
}
func readAll(filename string, notGetCalls bool) ([]Num_c, []error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Нет такого файла!")
		return []Num_c{}, []error{}
	}
	defer f.Close()
	// Чтение файла с ридером
	sc := bufio.NewScanner(f)
	Nums := []Num_c{}
	errs2 := []error{}
	for sc.Scan() {
		str := sc.Text()
		nn, errs := readline(str, notGetCalls)
		if errs == nil {
			Nums = append(Nums, nn.toNumC()...)
		} else {
			errs2 = append(errs2, errors.New("Ошибка "+errs.Error()+" в строке "+str))
		}

	}
	return Nums, errs2
}
func (n Num_c) getNum(ind int) int {
	l, err := strconv.Atoi(string(n.Num[ind]))
	if err != nil {
		log.Fatal("error Num conversion")
		return -1
	}
	return l
}
func (mt map_table) nullTable(k int) ([]int, []macrosline) {
	map_tbls := []int{}
	lines := []macrosline{}
	for _, line := range (*mt.Nums)[k] {
		if line.Transition {
			map_tblsN, linesN := mt.nullTable(line.Transit)
			map_tbls = append(map_tbls, map_tblsN...)
			lines = append(lines, linesN...)
		}
	}
	for i := range (*mt.Nums)[k] {
		lines = append(lines, macrosline{k, i, out_map{}})
	}
	map_tbls = append(map_tbls, k)
	(*mt.Nums)[k] = make([]out_map, 10)
	(*mt.Next) = append((*mt.Next), k)
	slices.Sort(*mt.Next)
	return map_tbls, lines
}

func transformline(n Num_c, mt *map_table, priority bool) (string, []int, []macrosline) {
	map_tabls := []int{}
	lines := []macrosline{}
	k := n.input
	if k >= len(*mt.Nums) {
		for i := len(*mt.Nums); i <= k; i++ {
			(*mt.Nums) = append((*mt.Nums), make([]out_map, 10))
			if i != k {
				(*mt.Next) = append((*mt.Next), i)
			}
		}
	}
	for i := 0; i < len(n.Num); i++ {
		nn := n.getNum(i)

		line := (*mt.Nums)[k][nn]
		if line.Transition {
			if priority {
				if i == len(n.Num)-1 {

					(*mt.Nums)[k][nn].Act_ptr = &n.act
					(*mt.Nums)[k][nn].Transition = false
					(*mt.Nums)[k][nn].Transit = 0
					map_tabls = append(map_tabls, k)
					lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
					map_tablsU, linesU := mt.upgrade(n)
					map_tabls = append(map_tabls, map_tablsU...)
					lines = append(lines, linesU...)

					map_tablsN, linesN := mt.nullTable(line.Transit)
					map_tabls = append(map_tabls, map_tablsN...)
					lines = append(lines, linesN...)
					return "", map_tabls, lines
				}
			} else {
				if i == len(n.Num)-1 {
					_, calls := getCalls([]Num_c{n}, mt.Nums)
					for _, act := range calls {
						if act.String() != ",000 - 0 - " {
							if act != n.act {
								return strconv.Itoa(n.input) + ":" + n.Num + ":" + n.act.String() + " - не может быть добавлен, имеет вызов: " + act.String(), map_tabls, lines
							}
						}
					}
					(*mt.Nums)[k][nn].Act_ptr = &n.act
					(*mt.Nums)[k][nn].Transition = false
					(*mt.Nums)[k][nn].Transit = 0
					map_tabls = append(map_tabls, k)
					lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
					map_tablsU, linesU := mt.upgrade(n)
					map_tabls = append(map_tabls, map_tablsU...)
					lines = append(lines, linesU...)

					map_tablsN, linesN := mt.nullTable(line.Transit)
					map_tabls = append(map_tabls, map_tablsN...)
					lines = append(lines, linesN...)
					return "", map_tabls, lines
				}
			}
			k = line.Transit
			continue
		}
		if line.Act_ptr != nil {
			if *line.Act_ptr == n.act {
				return "", map_tabls, lines
			} else {
				if priority {
					if i == len(n.Num)-1 {
						(*mt.Nums)[k][nn].Act_ptr = &n.act
						map_tabls = append(map_tabls, k)
						lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
						map_tablsU, linesU := mt.upgrade(n)
						map_tabls = append(map_tabls, map_tablsU...)
						lines = append(lines, linesU...)
						return "", map_tabls, lines
					} else {
						if len(*mt.Next) == 0 || (slices.IndexFunc((*mt.Next), func(c int) bool {
							return c > n.input
						}) == -1) {
							(*mt.Nums) = append((*mt.Nums), make([]out_map, 10))
							(*mt.Nums)[k][nn].Transition = true
							(*mt.Nums)[k][nn].Transit = len((*mt.Nums)) - 1
							(*mt.Nums)[k][nn].Act_ptr = nil
							map_tabls = append(map_tabls, k)
							lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
							//	(*mt.Next)[0] = len(*mt.Nums)

							k = len((*mt.Nums)) - 1
						} else {
							(*mt.Nums)[k][nn].Transition = true
							(*mt.Nums)[k][nn].Transit = (*mt.Next)[0]
							(*mt.Nums)[k][nn].Act_ptr = nil
							map_tabls = append(map_tabls, k)
							lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
							k = (*mt.Next)[0]
							*mt.Next = (*mt.Next)[1:]
						}
						for jj := i + 1; jj < len(n.Num); jj++ {
							nextn := n.getNum(jj)
							for j := 0; j <= 9; j++ {
								if nextn == j {
									continue
								}
								//	fmt.Println(n.Num[:(i+1)] + strconv.Itoa(j))
								_, mtt, ltt := transformline(Num_c{n.input, n.Num[:jj] + strconv.Itoa(j), *line.Act_ptr}, mt, false)
								map_tabls = append(map_tabls, mtt...)
								lines = append(lines, ltt...)
							}
						}

						continue

					}

				} else {
					return strconv.Itoa(n.input) + ":" + n.Num + ":" + n.act.String() + " - не может быть добавлен, имеет вызов: " + line.Act_ptr.String(), map_tabls, lines
				}
			}
		}
		if i == len(n.Num)-1 {
			(*mt.Nums)[k][nn].Act_ptr = &n.act
			map_tabls = append(map_tabls, k)
			lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
			map_tablsU, linesU := mt.upgrade(n)
			map_tabls = append(map_tabls, map_tablsU...)
			lines = append(lines, linesU...)
			return "", map_tabls, lines
		}
		if len(*mt.Next) == 0 || (slices.IndexFunc((*mt.Next), func(c int) bool {
			return c > n.input
		}) == -1) {
			(*mt.Nums) = append((*mt.Nums), make([]out_map, 10))
			(*mt.Nums)[k][nn].Transition = true
			(*mt.Nums)[k][nn].Transit = len((*mt.Nums)) - 1
			map_tabls = append(map_tabls, k)
			lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
			//	(*mt.Next)[0] = len(*mt.Nums)

			k = len((*mt.Nums)) - 1
		} else {
			(*mt.Nums)[k][nn].Transition = true
			(*mt.Nums)[k][nn].Transit = (*mt.Next)[0]
			map_tabls = append(map_tabls, k)
			lines = append(lines, macrosline{k, nn, (*mt.Nums)[k][nn]})
			k = (*mt.Next)[0]
			*mt.Next = (*mt.Next)[1:]
		}

	}
	return "", map_tabls, lines
}
func (mt *map_table) toMap(filename, macros string) {
	var f *os.File
	if mt.Nums != nil {
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			for i, arr := range *mt.Nums {
				io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
				for j, n := range arr {
					str := strconv.Itoa(j)
					if n.Transition {
						str += ": --> " + strconv.Itoa(n.Transit)
					} else if n.Act_ptr != nil {
						str += ": Вызов: " + (*n.Act_ptr).String()
					}
					io.WriteString(f, str)
					io.WriteString(f, "\r\n")
				}
				io.WriteString(f, "\r\n")
			}
			fmt.Println("Файл " + filename + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + filename + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(filename)
				f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				for i, arr := range *mt.Nums {
					io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
					for j, n := range arr {
						str := strconv.Itoa(j)
						if n.Transition {
							str += ": --> " + strconv.Itoa(n.Transit)
						} else if n.Act_ptr != nil {
							str += ": Вызов: " + (*n.Act_ptr).String()
						}
						io.WriteString(f, str)
						io.WriteString(f, "\r\n")
					}
					io.WriteString(f, "\r\n")

				}
				fmt.Println("Файл " + filename + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + filename + " не записан!")
			}
		}

		var ff *os.File
		if _, err := os.Stat(macros); errors.Is(err, os.ErrNotExist) {
			ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer ff.Close()
			for i, arr := range *mt.Nums {
				for j, n := range arr {
					if n.Transition {
						str := "wl map " + strconv.Itoa(i) + strconv.Itoa(j)
						str += " - - " + strconv.Itoa(n.Transit) + " 0 - -"
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					} else if n.Act_ptr != nil {
						str := "wl map " + strconv.Itoa(i) + strconv.Itoa(j)
						str += " c " + (*n.Act_ptr).String()
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					}
				}
			}
			fmt.Println("Файл " + macros + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + macros + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(macros)
				ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				defer ff.Close()
				for i, arr := range *mt.Nums {
					for j, n := range arr {
						if n.Transition {
							str := "wl map " + strconv.Itoa(i) + strconv.Itoa(j)
							str += " - - " + strconv.Itoa(n.Transit) + " 0 - -"
							io.WriteString(ff, str)
							io.WriteString(ff, "\r\n")
						} else if n.Act_ptr != nil {
							str := "wl map " + strconv.Itoa(i) + strconv.Itoa(j)
							str += " c " + (*n.Act_ptr).String()
							io.WriteString(ff, str)
							io.WriteString(ff, "\r\n")
						}
					}
				}
				fmt.Println("Файл " + macros + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + macros + " не записан!")
			}
		}
	} else {
		fmt.Println("База данных пуста")
	}
}
func (mt *map_table) save(filename string) {
	var err error
	json_data, err := json.Marshal(mt)
	if err != nil {
		fmt.Println("ошибка сохранения")
		return
	}
	var f *os.File
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		f, _ = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
		defer f.Close()
		io.WriteString(f, string(json_data))
		fmt.Println("Файл " + filename + " записан!")
		fmt.Println("База данных сохранена")
	} else {
		fmt.Println("Файл с именем " + filename + " уже существует. Заменить его? y/n")
		yn := ""
		fmt.Scanf("%s\r\n", &yn)
		if yn == "y" {
			os.Remove(filename)
			f, _ = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
			defer f.Close()
			io.WriteString(f, string(json_data))
			fmt.Println("Файл " + filename + " записан!")
			fmt.Println("База данных сохранена")
		} else {
			fmt.Println("Новый файл " + filename + " не записан!")
			fmt.Println("База данных не сохранена")
		}
	}
}

func load(filename string) (*map_table, error) {
	mt := map_table{}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	res := []byte{}
	for sc.Scan() {
		res = append(res, sc.Bytes()...)
	}
	err = json.Unmarshal(res, &mt)
	if err != nil {
		log.Fatal("error in unmarshall")
	}
	return &mt, nil
}

func transform(errs string, Nums []Num_c, slice_error *[]error) *map_table {
	mt := map_table{}
	mt.Nums = &([][]out_map{})
	mt.Next = &([]int{})
	for _, n := range Nums {
		err, _, _ := transformline(n, &mt, false)
		if err == "" {
			continue
		} else {
			*slice_error = append(*slice_error, errors.New(err))
		}

	}
	if len(*slice_error) != 0 {
		if _, err := os.Stat(errs); errors.Is(err, os.ErrNotExist) {
			fff, err1 := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
			if err1 != nil {
				panic(err)
			}
			defer fff.Close()
			for _, s := range *slice_error {
				io.WriteString(fff, s.Error()+"\r\n")
			}
			fmt.Println("Файл " + errs + " записан!")

			fmt.Println("Изменения внесены не полностью")
		} else {
			fmt.Println("Файл с именем " + errs + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(errs)
				fff, _ := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
				defer fff.Close()
				for _, s := range *slice_error {
					io.WriteString(fff, s.Error()+"\r\n")
				}
				fmt.Println("Файл " + errs + " записан!")

			} else {
				fmt.Println("Новый файл " + errs + " не записан!")
			}
			fmt.Println("База данных сформирована не полностью")
		}
	} else {
		fmt.Println("База данных сформирована успешно")
	}
	//	writeAll(&mt)
	return &mt
}
func (mt *map_table) addNew(filename, errs, macros string, Nums []Num_c, slice_error *[]error) {
	map_tbls := []int{}
	lines := []macrosline{}
	for _, n := range Nums {
		err, mtt, ltt := transformline(n, mt, false)
		map_tbls = append(map_tbls, mtt...)
		lines = append(lines, ltt...)
		if err == "" {
			continue
		} else {
			*slice_error = append(*slice_error, errors.New(err))
		}
	}
	var f *os.File
	if len(map_tbls) != 0 {
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			for i, arr := range *mt.Nums {
				if slices.Contains(map_tbls, i) {
					io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
					for j, n := range arr {
						str := strconv.Itoa(j)
						if n.Transition {
							str += ": --> " + strconv.Itoa(n.Transit)
						} else if n.Act_ptr != nil {
							str += ": Вызов: " + (*n.Act_ptr).String()
						}
						io.WriteString(f, str)
						io.WriteString(f, "\r\n")
					}
					io.WriteString(f, "\r\n")
				}
			}
			fmt.Println("Файл " + filename + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + filename + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(filename)
				f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				for i, arr := range *mt.Nums {
					if slices.Contains(map_tbls, i) {
						io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
						for j, n := range arr {
							str := strconv.Itoa(j)
							if n.Transition {
								str += ": --> " + strconv.Itoa(n.Transit)
							} else if n.Act_ptr != nil {
								str += ": Вызов: " + (*n.Act_ptr).String()
							}
							io.WriteString(f, str)
							io.WriteString(f, "\r\n")
						}
						io.WriteString(f, "\r\n")
					}
				}
				fmt.Println("Файл " + filename + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + filename + " не записан!")
			}
		}
	}
	if len(lines) != 0 {
		var ff *os.File
		if _, err := os.Stat(macros); errors.Is(err, os.ErrNotExist) {
			ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer ff.Close()
			for _, mline := range lines {

				if mline.line.Transition {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
					str += " - - " + strconv.Itoa(mline.line.Transit) + " 0 - -"
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				} else if mline.line.Act_ptr != nil {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
					str += " c " + (mline.line.Act_ptr).String()
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				} else {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn) + " - - - 0 - -"
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				}

			}
			fmt.Println("Файл " + macros + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + macros + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(macros)
				ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				defer ff.Close()
				for _, mline := range lines {

					if mline.line.Transition {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
						str += " - - " + strconv.Itoa(mline.line.Transit) + " 0 - -"
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					} else if mline.line.Act_ptr != nil {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
						str += " c " + (mline.line.Act_ptr).String()
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					} else {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn) + " - - - 0 - -"
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					}

				}
				fmt.Println("Файл " + macros + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + macros + " не записан!")
			}
		}
	}
	if len(*slice_error) != 0 {
		if _, err := os.Stat(errs); errors.Is(err, os.ErrNotExist) {
			fff, err1 := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
			if err1 != nil {
				panic(err)
			}
			defer fff.Close()
			for _, s := range *slice_error {
				io.WriteString(fff, s.Error()+"\r\n")
			}
			fmt.Println("Файл " + errs + " записан!")

			fmt.Println("Изменения внесены не полностью")
		} else {
			fmt.Println("Файл с именем " + errs + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(errs)
				fff, _ := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
				defer fff.Close()
				for _, s := range *slice_error {
					io.WriteString(fff, s.Error()+"\r\n")
				}
				fmt.Println("Файл " + errs + " записан!")

			} else {
				fmt.Println("Новый файл " + errs + " не записан!")
			}
		}
	}
}

func (mt *map_table) priorityAddNew(filename, errs, macros string, Nums []Num_c, slice_error *[]error) {
	map_tbls := []int{}
	lines := []macrosline{}
	for _, n := range Nums {
		err, mtt, ltt := transformline(n, mt, true)
		map_tbls = append(map_tbls, mtt...)
		lines = append(lines, ltt...)
		if err == "" {
			continue
		} else {
			*slice_error = append(*slice_error, errors.New(err))
		}
	}
	var f *os.File
	if len(map_tbls) != 0 {
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			for i, arr := range *mt.Nums {
				if slices.Contains(map_tbls, i) {
					io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
					for j, n := range arr {
						str := strconv.Itoa(j)
						if n.Transition {
							str += ": --> " + strconv.Itoa(n.Transit)
						} else if n.Act_ptr != nil {
							str += ": Вызов: " + (*n.Act_ptr).String()
						}
						io.WriteString(f, str)
						io.WriteString(f, "\r\n")
					}
					io.WriteString(f, "\r\n")
				}
			}
			fmt.Println("Файл " + filename + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + filename + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(filename)
				f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				for i, arr := range *mt.Nums {
					if slices.Contains(map_tbls, i) {
						io.WriteString(f, "########### "+strconv.Itoa(i)+" ###########\r\n")
						for j, n := range arr {
							str := strconv.Itoa(j)
							if n.Transition {
								str += ": --> " + strconv.Itoa(n.Transit)
							} else if n.Act_ptr != nil {
								str += ": Вызов: " + (*n.Act_ptr).String()
							}
							io.WriteString(f, str)
							io.WriteString(f, "\r\n")
						}
						io.WriteString(f, "\r\n")
					}
				}
				fmt.Println("Файл " + filename + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + filename + " не записан!")
			}
		}
	}
	if len(lines) != 0 {
		var ff *os.File
		if _, err := os.Stat(macros); errors.Is(err, os.ErrNotExist) {
			ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer ff.Close()
			for _, mline := range lines {

				if mline.line.Transition {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
					str += " - - " + strconv.Itoa(mline.line.Transit) + " 0 - -"
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				} else if mline.line.Act_ptr != nil {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
					str += " c " + (mline.line.Act_ptr).String()
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				} else {
					str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn) + " - - - 0 - -"
					io.WriteString(ff, str)
					io.WriteString(ff, "\r\n")
				}

			}
			fmt.Println("Файл " + macros + " успешно записан!")
		} else {
			fmt.Println("Файл с именем " + macros + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(macros)
				ff, err = os.OpenFile(macros, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
				if err != nil {
					log.Fatal(err)
				}
				defer ff.Close()
				for _, mline := range lines {

					if mline.line.Transition {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
						str += " - - " + strconv.Itoa(mline.line.Transit) + " 0 - -"
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					} else if mline.line.Act_ptr != nil {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn)
						str += " c " + (mline.line.Act_ptr).String()
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					} else {
						str := "wl map " + strconv.Itoa(mline.n) + strconv.Itoa(mline.nn) + " - - - 0 - -"
						io.WriteString(ff, str)
						io.WriteString(ff, "\r\n")
					}

				}
				fmt.Println("Файл " + macros + " успешно записан!")
			} else {
				fmt.Println("Новый файл " + macros + " не записан!")
			}
		}
	}
	if len(*slice_error) != 0 {
		if _, err := os.Stat(errs); errors.Is(err, os.ErrNotExist) {
			fff, err1 := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
			if err1 != nil {
				panic(err)
			}
			defer fff.Close()
			for _, s := range *slice_error {
				io.WriteString(fff, s.Error()+"\r\n")
			}
			fmt.Println("Файл " + errs + " записан!")

			fmt.Println("Изменения внесены не полностью")
		} else {
			fmt.Println("Файл с именем " + errs + " уже существует. Заменить его? y/n")
			yn := ""
			fmt.Scanf("%s\r\n", &yn)
			if yn == "y" {
				os.Remove(errs)
				fff, _ := os.OpenFile(errs, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
				defer fff.Close()
				for _, s := range *slice_error {
					io.WriteString(fff, s.Error()+"\r\n")
				}
				fmt.Println("Файл " + errs + " записан!")

			} else {
				fmt.Println("Новый файл " + errs + " не записан!")
			}
		}
	}
}

func main() {
	command := ""
	file := ""
	file1 := ""
	file2 := ""
	file3 := ""
	mt := new(map_table)
	fmt.Println("Программа готова к работе")
	for command != "exit" {
		fmt.Scanf("%s %s %s %s %s\r\n", &command, &file, &file1, &file2, &file3)
		switch command {
		case "new":
			{
				if len(file) < 2 {
					fmt.Println("Неправильное имя файла")
					fmt.Println("Команда не выполнена, help - для справки")
					break
				} else {
					if string(file[0]) != "\"" {
						fmt.Println("Неправильное имя файла")
						fmt.Println("Команда не выполнена, help - для справки")
						break
					} else {
						file = file[1 : len(file)-1]
					}
				}
				if len(file1) < 2 {
					file1 = "errs.txt"
				} else {
					if string(file1[0]) != "\"" {
						file1 = "errs.txt"
					} else {
						file1 = file1[1 : len(file1)-1]
					}
				}
				n, err := readAll(file, true)
				mt = transform(file1, n, &err)

			}
		case "load":
			{
				if len(file) < 2 {
					fmt.Println("Неправильное имя файла")
					fmt.Println("Команда не выполнена, help - для справки")
					break
				} else {
					if string(file[0]) != "\"" {
						fmt.Println("Неправильное имя файла")
						fmt.Println("Команда не выполнена, help - для справки")
						break
					} else {
						file = file[1 : len(file)-1]
					}
				}
				var err error
				mt, err = load(file)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("База данных загружена")
				}

			}
		case "save":
			{
				if len(file) < 2 {
					file = "db.json"
				} else {
					if string(file[0]) != "\"" {
						file = "db.json"
					} else {
						file = file[1 : len(file)-1]
					}
				}
				mt.save(file)
			}
		case "getCalls":
			{
				if len(file) < 2 {
					fmt.Println("Неправильное имя файла")
					fmt.Println("Команда не выполнена, help - для справки")
					break
				} else {
					if string(file[0]) != "\"" {
						fmt.Println("Неправильное имя файла")
						fmt.Println("Команда не выполнена, help - для справки")
						break
					} else {
						file = file[1 : len(file)-1]
					}
				}
				if len(file1) < 2 {
					file1 = "calls.txt"
				} else {
					if string(file1[0]) != "\"" {
						file1 = "calls.txt"
					} else {
						file1 = file1[1 : len(file1)-1]
					}
				}
				n, err := readAll(file, false)
				if len(err) > 0 {
					for _, er := range err {
						fmt.Println(er)
					}
				} else {
					nn, a := getCalls(n, mt.Nums)
					printCalls(file1, nn, a)

				}

			}
		case "result":
			{

				if len(file) < 2 {
					file = "map.txt"
				} else {
					if string(file[0]) != "\"" {
						file = "map.txt"
					} else {
						file = file[1 : len(file)-1]
					}
				}
				if len(file1) < 2 {
					file1 = "macros.txt"
				} else {
					if string(file1[0]) != "\"" {
						file1 = "macros.txt"
					} else {
						file1 = file1[1 : len(file1)-1]
					}
				}
				mt.toMap(file, file1)

			}
		case "add":
			{
				if len(file) < 2 {
					fmt.Println("Неправильное имя файла")
					fmt.Println("Команда не выполнена, help - для справки")
					break
				} else {
					if string(file[0]) != "\"" {
						fmt.Println("Неправильное имя файла")
						fmt.Println("Команда не выполнена, help - для справки")
						break
					} else {
						file = file[1 : len(file)-1]
					}
				}
				if len(file1) < 2 {
					file1 = "errs.txt"
				} else {
					if string(file1[0]) != "\"" {
						file1 = "errs.txt"
					} else {
						file1 = file1[1 : len(file)-1]
					}
				}
				if len(file2) < 2 {
					file2 = "map.txt"
				} else {
					if string(file2[0]) != "\"" {
						file2 = "map.txt"
					} else {
						file2 = file2[1 : len(file2)-1]
					}
				}
				if len(file3) < 2 {
					file3 = "macros.txt"
				} else {
					if string(file3[0]) != "\"" {
						file3 = "macros.txt"
					} else {
						file3 = file3[1 : len(file3)-1]
					}
				}
				n, err := readAll(file, true)
				mt.addNew(file2, file1, file3, n, &err)
			}
		case "priorityAdd":
			{
				if len(file) < 2 {
					fmt.Println("Неправильное имя файла")
					fmt.Println("Команда не выполнена, help - для справки")
					break
				} else {
					if string(file[0]) != "\"" {
						fmt.Println("Неправильное имя файла")
						fmt.Println("Команда не выполнена, help - для справки")
						break
					} else {
						file = file[1 : len(file)-1]
					}
				}
				if len(file1) < 2 {
					file1 = "errs.txt"
				} else {
					if string(file1[0]) != "\"" {
						file1 = "errs.txt"
					} else {
						file1 = file1[1 : len(file1)-1]
					}
				}
				if len(file2) < 2 {
					file2 = "map.txt"
				} else {
					if string(file2[0]) != "\"" {
						file2 = "map.txt"
					} else {
						file2 = file2[1 : len(file2)-1]
					}
				}
				if len(file3) < 2 {
					file3 = "macros.txt"
				} else {
					if string(file3[0]) != "\"" {
						file3 = "macros.txt"
					} else {
						file3 = file3[1 : len(file3)-1]
					}
				}
				n, err := readAll(file, true)
				mt.priorityAddNew(file2, file1, file3, n, &err)
			}
		case "help":
			{
				fmt.Println("new \"путь_к_файлу_input\" [\"путь_к_файлу_errs\"]")
				fmt.Println("load \"путь_к_файлу_с_БД\"")
				fmt.Println("save [\"путь_к_файлу_с_БД\"]")
				fmt.Println("result [\"путь_к_файлу_с_map\"] [\"путь_к_файлу_macros\"]")
				fmt.Println("getCalls \"путь_к_файлу_input_calls\" [\"путь_к_файлу_errs\"]")
				fmt.Println("add \"путь_к_файлу_input_add\" [\"путь_к_файлу_errs\"] [\"путь_к_файлу_macros\"]")
				fmt.Println("priorityAdd \"путь_к_файлу_input_add\" [\"путь_к_файлу_errs\"] [\"путь_к_файлу_macros\"]")
				fmt.Println("help")
				fmt.Println("Подробнее смотрите в файле help")
			}
		case "exit":
			{
				fmt.Println("Выход...")
			}
		default:
			{
				fmt.Println(command + " - нет такой команды, help - для справки")
			}
		}
	}
	/*
		start_Nums := readAll("./DX/input")
		m := transform(start_Nums)
		//	m := load("json")
		//fmt.Println((*m.Nums)[4])
		n := readline("0:00000-99999:gr,02").toNumC()
		printCalls(getCalls(n, m.Nums))
		/*
			n := transform(start_Nums)
			nn := readline("10:62113-62120:Gr,02").toNumC()
			printCalls(getCalls(nn, n.Nums), nn)*/

	/*
	   n := readline("0:62110:Gr,02").toNumC()
	   a, err := n[0].getLastMap(&Nums)

	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   fmt.Println(a)
	*/

}
