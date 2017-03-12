package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func start() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", handler)
	http.HandleFunc("/cupcakes", cupcakes)
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})
	http.Handle("/assets/", http.FileServer(http.Dir("./")))
}

func cupcakes(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("./assets/cupcake.jpeg")
	w.Write(data)
}

func handler(w http.ResponseWriter, r *http.Request) {
	view, _ := template.ParseFiles("./views/index.html")
	view.Execute(w, nil)
}

func handleWebSocket(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.register <- conn
	for {
		_, msg, err := conn.ReadMessage()
		fmt.Println(string(msg))
		if err != nil {
			fmt.Println(err)
			h.unregister <- conn
			return
		}
		h.broadcast <- msg
	}
}

func main() {
	start()
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
