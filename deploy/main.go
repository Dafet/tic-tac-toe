package main

import (
	"tic-tac-toe/matchmaking"
	"tic-tac-toe/ws"
	"time"
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

	// go func() {
	// 	for {
	// 		fmt.Println("broadcasting...")
	// 		server.BroadcastMsg(ws.NewSetUserDataMsg("sample"))
	// 		time.Sleep(time.Millisecond * 500)
	// 	}
	// }()

	mm := matchmaking.NewRandom(time.Duration(time.Millisecond * 500))
	server.SetMatchmakingEngine(mm)

	server.Start()

	// fmt.Scanln()
}
