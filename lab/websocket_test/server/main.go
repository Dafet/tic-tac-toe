package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	_ = StartServer(messageHandler)

	// make as separate test func if needed
	// for {
	// 	msg := "Hello"
	// 	fmt.Printf("writing msg: %s \n", msg)
	// 	server.WriteMessage([]byte(msg))
	// 	time.Sleep(time.Second * 1)
	// }
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte) // хандлер новых сообщений
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(":8080", nil) // Уводим http сервер в горутину

	fmt.Println("started listen on :8080")

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {

	// panicking here if request has been done via browser

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errMsg := "error upgrading connection: " + err.Error()
		fmt.Println(errMsg)
		fmt.Fprintln(w, errMsg)
	}

	defer connection.Close()

	server.clients[connection] = true        // Сохраняем соединение, используя его как ключ
	defer delete(server.clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		go server.handleMessage(message)
	}
}

func (server *Server) WriteMessage(message []byte) {
	for conn := range server.clients {
		time.Sleep(time.Second * 1)
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
