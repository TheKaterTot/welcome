package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func start() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/cupcakes", cupcakes)
}

func cupcakes(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("./cupcake.jpeg")
	w.Write(data)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome, fans of Katy")
}

func main() {
	start()
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
