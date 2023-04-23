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
}
func main() {
	c := make(chan string) // create the chanel
	go count("sheep", c)   // create the thread and pass the channel
	msg := <-c             // get a message from the channel
	fmt.Println(msg)       // print it out (will only print one message)
}
