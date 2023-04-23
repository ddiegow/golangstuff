package main

// TODO: rewrite to use worker-coordinator system using channels
// TODO: figure out that found: https://golang.org/cmd/ "" output
import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}
type SafePool struct {
	mu    sync.Mutex
	pool  map[string]int
	found map[string]int
}

func worker(safePool *SafePool, fetcher Fetcher) { // first thread
	for { // infinite loop

		safePool.mu.Lock()
		if len(safePool.pool) == 0 {
			safePool.mu.Unlock()
			break
		}

		var theUrl string
		for k := range safePool.pool {
			theUrl = k // get the first key
			break
		}
		theDepth := safePool.pool[theUrl]
		delete(safePool.pool, theUrl)

		if theDepth > 0 && len(theUrl) > 0 {
			body, urls, err := fetcher.Fetch(theUrl)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("found: %s %q\n", theUrl, body)
			safePool.found[theUrl] = 1
			for _, u := range urls {
				safePool.pool[u] = theDepth - 1
			}
		}
		safePool.mu.Unlock()

	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	fmt.Println("Starting crawler")
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	// create a pool that will hold all the urls that need processing
	safePool := SafePool{pool: make(map[string]int), found: make(map[string]int)}
	safePool.pool[url] = depth
	go worker(&safePool, fetcher)
	go worker(&safePool, fetcher)
	fmt.Scanln()
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

/*
EXPECTED RESULTS
found: https://golang.org/ "The Go Programming Language"
found: https://golang.org/pkg/ "Packages"
not found: https://golang.org/cmd/
found: https://golang.org/pkg/fmt/ "Package fmt"
found: https://golang.org/pkg/os/ "Package os"
*/
