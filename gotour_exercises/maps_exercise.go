package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	separated := strings.Fields(s)
	result := make(map[string]int)
	for _, v := range separated {
		if result[v] == 0 {
			result[v] = 1
		} else {
			result[v]++
		}
	}

	return result
}

func main() {
	wc.Test(WordCount)
}
