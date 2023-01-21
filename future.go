package main

import (
	"fmt"
	"net/http"
	"time"
)

////////////////////
// Simple futures

// A future, once available, will be transmitted via a channel.
// The Boolean parameter indicates if the (future) computation succeeded or failed.

type Comp struct {
	val    interface{}
	status bool
}

type Future chan Comp

type Futures []Future

func (ft Future) onSuccess(cb func(interface{})) {
	go func() {
		v, o := ft.get()
		if o {
			cb(v)
		}
	}()
}

func (ft Future) onFailure(cb func()) {
	go func() {
		_, o := ft.get()
		if !o {
			cb()
		}
	}()
}

func (ft Future) then(cb func() (interface{}, bool)) Future {
	return future(cb)
}

func (f Future) get() (interface{}, bool) {
	v := <-f
	return v.val, v.status
}

func future(f func() (interface{}, bool)) Future {
	ch := make(chan Comp)
	go func() {
		r, s := f()
		v := Comp{r, s}
		for {
			ch <- v
		}
	}()
	return ch
}

func (fts Futures) mapsOn(cb func(interface{})) {
	for _, el := range fts {
		go func(ft Future) {
			v, o := ft.get()
			if o {
				cb(v)
			} else {

			}
		}(el)
	}
}

func (fts Futures) getAll() []interface{} {
	output := make([]interface{}, len(fts))
	for i, el := range fts {
		v, o := el.get()
		if o {
			output[i] = v
		} else {
			fmt.Print("false")
		}
	}
	return output
}

///////////////////////////////
// Adding more functionality

// Pick first available future
func (ft Future) first(ft2 Future) Future {

	return future(func() (interface{}, bool) {

		var v interface{}
		var o bool

		// check for any result to become available
		select {
		case x := <-ft:
			v = x.val
			o = x.status

		case x2 := <-ft2:
			v = x2.val
			o = x2.status

		}

		return v, o
	})
}


func (fts Futures) first() Future {
	return future(func() (interface{}, bool) {
		var v interface{}
		var o bool
		agg := make(chan Future)
		for _, ch := range fts {
			go func(ft Future) {
				agg <- ft
			}(ch)
		}

		select {
		case ft := <-agg:
			x := <-ft
			v = x.val
			o = x.status
		}
		return v, o
	})
}

// Pick first successful future
func (ft Future) firstSucc(ft2 Future) Future {

	return future(func() (interface{}, bool) {

		var v interface{}
		var o bool

		select {
		case x := <-ft:
			if x.status {
				v = x.val
				o = x.status
			} else {
				v, o = ft2.get()
			}

		case x2 := <-ft2:
			if x2.status {
				v = x2.val
				o = x2.status
			} else {
				v, o = ft.get()
			}

		}

		return v, o
	})
}

///////////////////////
// Examples

func getSite(url string) Future {
	return future(func() (interface{}, bool) {
		resp, err := http.Get(url)
		if err == nil {
			return resp, true
		}
		return err, false
	})
}

func printResponse(response *http.Response) {
	fmt.Println(response.Request.URL)
	header := response.Header
	// fmt.Println(header)
	date := header.Get("Date")
	fmt.Println(date)

}

func example1() {

	stern := getSite("http://www.stern.de")

	stern.onSuccess(func(result interface{}) {
		response := result.(*http.Response)
		printResponse(response)

	})

	stern.onFailure(func() {
		fmt.Printf("failure \n")
	})

	fmt.Printf("do something else \n")

	time.Sleep(2 * time.Second)

}

func example2() {

	spiegel := getSite("http://www.spiegel.de")
	stern := getSite("http://www.stern.de")
	welt := getSite("http://www.welt.com")

	req := spiegel.first(stern.first(welt))

	req.onSuccess(func(result interface{}) {
		response := result.(*http.Response)
		printResponse(response)

	})

	req.onFailure(func() {
		fmt.Printf("failure \n")
	})

	fmt.Printf("do something else \n")

	time.Sleep(2 * time.Second)

}

func main() {

	spiegel := getSite("http://www.spiegel.de")
	stern := getSite("http://www.stern.de")
	welt := getSite("http://www.welt.com")

	fts := Futures{spiegel, stern, welt}
	/*
		fts.mapsOn(func(result interface{}) {
			response := result.(*http.Response)
			printResponse(response)
		})*/
	fut := fts.first()
	fmt.Printf("do something else \n")
	v, _ := fut.get()
	printResponse(v.(*http.Response))
	time.Sleep(2 * time.Second)
}
