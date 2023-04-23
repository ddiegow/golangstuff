package main

import (
	"fmt"
	"time"
)

func count(s string) {
	for i := 0; true; i++ {
		fmt.Println(i, s)
		time.Sleep(time.Millisecond * 500)
	}
}
func main() {
	go count("sheep")
	count("fish")
}
