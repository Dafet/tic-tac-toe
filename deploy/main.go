package main

import (
	"tic-tac-toe/ws"
)

// inject logger everywhere?

func main() {
	server := ws.New()

	// go func() {
	// 	for {
	// 		server.LogConns()
	// 		time.Sleep(time.Millisecond * 500)
	// 	}
	// }()

	server.Start()

	// for {
	// 	fmt.Println("broadcasting: Hello client")
	// 	server.BroadcastMsg("Hello client")
	// 	time.Sleep(time.Millisecond * 500)
	// }

	// fmt.Scanln()
}
