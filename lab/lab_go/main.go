package main

import (
	"fmt"
	"time"
)

func main() {
	sample()

	time.Sleep(time.Second * 10)
}

func sample() {
	go func() {
		for {
			time.Sleep(time.Second * 1)

			fmt.Println("hello from go")
		}
	}()
}
