package main

import (
	"fmt"
	"time"
)

// Modeling wait/notify in Go with channels

type Group chan int

func newGroup() Group {
	return make(chan int)
}

// Wait till notified
func (g Group) wait() {
	<-g
}

func (g Group) waitSomeTime(s time.Duration) {
	select {
	case <-g:
	case <-time.After(s):
	}

}

// Notify one of the waiting threads.
// If nobody is waiting, the signal gets lost.
func (g Group) notify() {
	select {
	case g <- 1:
	default:
	}
}

// notifyAll
// Loop till all waiting threads are notified.
func (g Group) notifyAll() {
	b := true
	for b {
		select {
		case g <- 1:
		default:
			b = false
		}
	}
}

// Sleeping barber example making use of wait/notify
func sleepingBarber() {
	g := newGroup()

	customer := func(s string) {
		for {
			g.wait()
			fmt.Printf("%s got haircut! \n", s)
			time.Sleep(1 * time.Second)

		}
	}

	barber := func() {
		for {
			g.notify() // single barber checks for waiting customer
			// g.notifyAll()  // as many barbers as there are waiting customers
			fmt.Printf("cut hair! \n")
			time.Sleep(3 * time.Second)
		}

	}

	go customer("A")
	go customer("B")
	go customer("C")

	barber()

}

func main() {

	sleepingBarber()

}
