package main

import (
	"fmt"
	"time"
)

// MVar mit Hilfe
type MVar (chan int)

func newMVar(x int) MVar {
	var ch = make(chan int)
	go func() { ch <- x }()
	return ch
}

func takeMVar(m MVar) int {
	var x int
	x = <-m
	return x
}

func putMVar(m MVar, x int) {
	m <- x
}

// MVar Beispiel
func producer(m MVar) {
	var x int = 1
	for {
		time.Sleep(1 * 1e9)
		putMVar(m, x)
		x++
	}
}

func consumer(m MVar) {
	for {
		var x int = takeMVar(m)
		fmt.Printf("Received %d \n", x)
	}
}

func testMVar() {
	var m MVar

	m = newMVar(1)

	go producer(m)

	consumer(m)

}

// MVar Beispiel 2
// 2 ueberholt 1
func testMVar2() {
	m := newMVar(1)  // 1
	go putMVar(m, 2) // 2
	x := takeMVar(m)
	fmt.Printf("Received %d \n", x)
}

// Deadlock
func testMVar3() {
	var m MVar
	m = newMVar(1) // Full
	takeMVar(m)    // Empty
	putMVar(m, 2)  // Full
}

// 2te MVar Kodierung
const (
	Empty = 0
	Full  = 1
)

func newMVar2(x int) MVar {
	var ch = make(chan int)
	go func() {
		var state = Full
		var elem int = x
		for {
			switch {
			case state == Full:
				ch <- elem
				state = Empty
			case state == Empty:
				elem = <-ch
				state = Full
			}
		}
	}()
	return ch
}

// Wir verwenden newMVar2 anstatt newMVar
func testMVar4() {
	m := newMVar2(1) // 1
	go putMVar(m, 2) // 2
	x := takeMVar(m)
	fmt.Printf("Received %d \n", x)
}

func testMVar5() {
	var m MVar
	m = newMVar2(1) // Full
	takeMVar(m)     // Empty
	putMVar(m, 2)   // Full
}

func main() {

	// testMVar()
	// testMVar2()
	// testMVar3()

	testMVar4()
	testMVar5()

}
