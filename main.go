package main

func main() {
	x := make(chan int)
	var v int
	go func() { v = <-x }()
	print(v)
}
