package main

// Improved with inspiration from Harvard 6.824 lecture
import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
}
type SafePool struct {
	mu      sync.Mutex
	fetched map[string]bool
}

// The only time we need to lock and unlock is when working on the state that is shared by all threads
// Putting the routine in a method is cleaner than putting it in the Crawl function, though it would work as well
func (safePool *SafePool) processState(url string) bool {
	safePool.mu.Lock()               // look the pool so there are no conflicts
	defer safePool.mu.Unlock()       // unlock when we exit the function
	fetched := safePool.fetched[url] // get the "fetched" state of the url
	safePool.fetched[url] = true     // set the "fetched" state to true
	return fetched                   // return the original "fetched" state
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// We skip the depth constraint, just scan everything
func Crawl(url string, fetcher Fetcher, safePool *SafePool) {
	if safePool.processState(url) { // if we already fetched it
		return // do nothing
	}
	urls, err := fetcher.Fetch(url) // fetch the urls
	if err != nil {                 // if there was an error
		return // do nothing
	}
	var wg sync.WaitGroup    // will use to wait for other goroutines to finish. if we don't use it, the program would only run the Crawl function once and exit
	for _, u := range urls { // process the urls
		wg.Add(1)           // add a task
		go func(u string) { // inline function takes an argument so that all goroutines don't share the same u
			defer wg.Done()             // when the crawl finishes, we're done
			Crawl(u, fetcher, safePool) // call Crawl recursively with one less level of depth
		}(u)
	}
	wg.Wait() // wait for other routines to finish
}

func main() {
	safePool := SafePool{fetched: make(map[string]bool)} // create a safepool for all goroutines to share
	Crawl("https://golang.org/", fetcher, &safePool)     // start the crawl
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found: %s %q\n", url, res.body) // print it here. could do it in the crawler body, but this is cleaner
		return res.urls, nil
	}
	fmt.Printf("not found: %s\n", url) // same as previous comment
	return nil, fmt.Errorf("not found: %s", url)
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
