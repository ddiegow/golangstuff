package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	fib := -1
	pre := 0
	return func() int {
		if fib == -1 { // if we've just started
			fib = 0 // return zero
		} else if fib == 0 { // if we're at zero
			pre = 0 // store current value
			fib = 1 // return next fibonacci value
		} else { // if we've passed zero
			tmp := pre      // store the previous value in tmp
			pre = fib       // store the current fibonacci value
			fib = fib + tmp // add the previous fibonacci value to the current one to obtain the next fibonacci value
		}
		return fib // return the value
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
