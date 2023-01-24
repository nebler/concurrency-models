package main

import (
	"fmt"
	"math/rand"
	"time"
)

////////////////////
// Reducing the number of gourtines.
// Internally, we use a promise as we might want to explicitly provide the (future) value.

// A promise is a form of a future where the result can be provided explicitely via some setSucc/setFail primitives.
// A promise can only be completed once.
// We keep of the completion status via the boolean flag empy.
// If empty any callback will be registered and executed once the promise is completed.

type Promise[T any] struct {
	val           T
	status        bool
	m             chan int
	succCallBacks []func(T)
	failCallBacks []func()
	empty         bool
}

func newPromise[T any]() *Promise[T] {
	p := Promise[T]{empty: true, m: make(chan int, 1), succCallBacks: make([]func(T), 0), failCallBacks: make([]func(), 0)}
	return &p
}

// Check if still empty.
// If yes, call all registered callbacks within one goroutine.
// Otherwise, do nothing.
func (p *Promise[T]) setSucc(v T) {
	p.m <- 1
	if p.empty {
		p.val = v
		p.status = true
		p.empty = false
		succs := p.succCallBacks
		p.succCallBacks = make([]func(T), 0)
		<-p.m
		go func() {
			for _, cb := range succs {
				cb(v)
			}
		}()
	} else {
		<-p.m
	}

}

func (p *Promise[T]) setFail() {
	p.m <- 1
	if p.empty {
		p.status = false
		p.empty = false
		fails := p.failCallBacks
		p.failCallBacks = make([]func(), 0)
		<-p.m
		go func() {
			for _, cb := range fails {
				cb()
			}
		}()
	} else {
		<-p.m
	}

}

func future[T any](f func() (T, bool)) *Promise[T] {
	p := newPromise[T]()
	go func() {
		r, s := f()
		if s {
			p.setSucc(r)
		} else {
			p.setFail()
		}
	}()
	return p
}

func (p *Promise[T]) complete(f func() (T, bool)) {
	go func() {
		r, s := f()
		if s {
			p.setSucc(r)
		} else {
			p.setFail()
		}
	}()

}

func (p *Promise[T]) onSuccess(cb func(T)) {
	p.m <- 1
	switch {
	case p.empty:
		p.succCallBacks = append(p.succCallBacks, cb)
	case !p.empty && p.status:
		go cb(p.val)
	default: // drop cb, will never be called

	}
	<-p.m

}

func (p *Promise[T]) onFailure(cb func()) {
	p.m <- 1
	switch {
	case p.empty:
		p.failCallBacks = append(p.failCallBacks, cb)
	case !p.empty && !p.status:
		go cb()
	default: // drop cb, will never be called

	}
	<-p.m

}

///////////////////////////////
// Adding more functionality

// Try to complete p with p2.
// We only consider the successful case.
func (p *Promise[T]) tryCompleteWith(p2 *Promise[T]) {
	p2.onSuccess(func(v T) {
		p.setSucc(v)
	})

}

// Pick first successful future
func (p *Promise[T]) firstSucc(p2 *Promise[T]) *Promise[T] {
	p3 := newPromise[T]()
	p3.tryCompleteWith(p)
	p3.tryCompleteWith(p2)
	return p3
}

///////////////////////
// Examples

// Holiday booking
func example1() {

	// Book some Hotel. Report price (int) and some poential failure (bool).
	booking := func() (int, bool) {
		// time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
		return rand.Intn(50), true
	}

	f1 := newPromise[int]()
	f1.complete(booking)

	f2 := newPromise[int]()
	f2.complete(booking)

	f3 := f1.firstSucc(f2)

	f3.onSuccess(func(quote int) {

		fmt.Printf("\n Hotel asks for %d Euros", quote)
	})

	time.Sleep(2 * time.Second)
}

func example2() {

	// Book some Hotel. Report price (int) and some poential failure (bool).
	booking := func() (int, bool) {
		// time.Sleep((time.Duration)(rand.Intn(999)) * time.Millisecond)
		return rand.Intn(50), true
	}

	f1 := future[int](booking)

	f2 := future[int](booking)

	f3 := f1.firstSucc(f2)

	f3.onSuccess(func(quote int) {

		fmt.Printf("\n Hotel asks for %d Euros", quote)
	})

	time.Sleep(2 * time.Second)
}

func ex3() {

}

func main() {

	example2()
}
