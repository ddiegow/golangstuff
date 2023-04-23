package main

import (
	"fmt"
	"time"
)

func count(s string, c chan string) {
	for i := 0; i < 6; i++ {
		c <- s // send the string to the message
		time.Sleep(time.Millisecond * 500)
	}
	close(c) // if we don't close this, the main thread will be waiting forever for a thread to send data through the channel
}
func main() {
	c := make(chan string) // create the chanel
	go count("sheep", c)   // create the thread and pass a channel
	/*
		for {                  // infinite loop
			msg, open := <-c // get both the data AND the open status value
			if !open {       // if the channel is closed
				break // exit loop
			}
			fmt.Println(msg) // print out what's in the channel
		}
	*/
	// Even better with a for range, that way we don't need to check the value of open (it's done automatically)
	for msg := range c {
		fmt.Println(msg)
	}
	// close(c) NEVER DO THIS, RECEIVERS MUST NEVER EVER CLOSE CHANNELS

	// we can also create a buffered receiver
	c = make(chan string, 2)
	c <- "hello"
	c <- "world"
	// c <- "oops" this would overflow the buffer and the thread would block until the input is read, which would block the main thread
	fmt.Println(<-c)
	fmt.Println(<-c)
}
