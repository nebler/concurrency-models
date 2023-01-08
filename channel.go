package main

import (
	"fmt"
	"time"
)

func snd(s string, ch chan int) {
	var x int = 0
	for {
		x++
		ch <- x
		fmt.Printf("%s sendet %d \n", s, x)
		time.Sleep(1 * 1e9)
	}

}

func rcv(ch chan int) {
	var x int
	for {
		x = <-ch
		fmt.Printf("empfangen %d \n", x)

	}

}

func main() {
	var ch chan int = make(chan int)
	go snd("A", ch)
	go rcv(ch)
	for {

	}
}
