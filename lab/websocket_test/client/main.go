// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// (return?)go:build ignore
// (return?)+build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"tic-tac-toe/server/ws"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// infiniteWriteReadLoop(c) // test

	startTTTGame(c)
}

func startTTTGame(c *websocket.Conn) {
	fmt.Println("type player's name: ")

	var name string

	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatalln("[error] reading player's name: ", err)
	}

	msg, err := ws.NewPlayRdyMsgSerialized(name)
	if err != nil {
		log.Fatalln("[error] creating rdy msg: ", err)
	}

	// make msg send utils

}

func infiniteWriteReadLoop(c *websocket.Conn) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			// writeMsgSample(c, t) // working sample
			writeTestMsg(c)
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func writeMsgSample(c *websocket.Conn, t time.Time) {
	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

// todo: rewrite/delete
func writeTestMsg(c *websocket.Conn) {
	tm := ws.NewTestMsg()

	data, err := tm.Serialize()
	if err != nil {
		log.Println("[error] serializing msg: ", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("[error] error writing msg: ", err)
		return
	}
}
