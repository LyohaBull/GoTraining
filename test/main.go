package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"
)

func test(c chan string) {
	a := <-c
	fmt.Println(a)
	c <- "sad"
	b := <-c
	fmt.Println(b)
}
func testR(w http.ResponseWriter, r *http.Request) {
	//http.Handle("/src/", http.FileServer(http.Dir("/src/")))
	tmp, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	tmp.Execute(w, nil)
}

type data struct {
	messages []string
}

func z() func(http.ResponseWriter, *http.Request) {
	var d = data{}
	return func(w http.ResponseWriter, r *http.Request) {
		mess, _ := io.ReadAll(r.Body)
		d.messages = append(d.messages, string(mess))
		fmt.Println(d.messages)
		w.Write([]byte(fmt.Sprintf("response %d", len(d.messages))))
		if len(d.messages) > 1 {
			d.messages = []string{}
		}

	}
}

func testA(w http.ResponseWriter, r *http.Request) {
	//http.Handle("/src/", http.FileServer(http.Dir("/src/")))

}
func main() {
	//
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))
	http.HandleFunc("/", testR)
	http.HandleFunc("/auth", z())
	http.ListenAndServe(":8080", nil)
}
