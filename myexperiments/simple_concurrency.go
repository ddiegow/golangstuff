package main

import (
	"fmt"
	"sync"
	"time"
)

func count(s string) {
	for i := 0; i < 6; i++ {
		fmt.Println(i, s)
		time.Sleep(time.Millisecond * 500)
	}
}
func main() {
	var wg sync.WaitGroup // when the variable isn't initialized, it needs to be created with var
	wg.Add(1)             // we add to the counter of running threads
	go func() {
		count("sheep") // call the function
		wg.Done()      // when it's done, we call WaitGroup.Done() to decrease the thread counter
	}()
	wg.Wait() // wait until the thread counter is zero

}
