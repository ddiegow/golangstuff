package main

// TODO: someday make a version that doesn't use a shared pool between the goroutines, something using channels.
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
	mu        sync.Mutex
	processed map[string]int
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, safePool *SafePool) {
	// first check depth
	if depth <= 0 { // if lower or equal to zero, send back nil
		return
	}
	body, urls, err := fetcher.Fetch(url) // fetch the urls
	if err != nil {                       // if there was an error
		fmt.Println(err)                // print it out
		safePool.processed[url] = depth // add the url to the processed list
		return                          // we're done
	}
	fmt.Printf("found: %s %q\n", url, body) // no error, so print out the url
	safePool.processed[url] = depth         // add it to the processed list

	var wg sync.WaitGroup    // will use to wait for other goroutines to finish. if we don't use it, the program would only run the Crawl function once and exit
	safePool.mu.Lock()       // lock the pool so there are no conflicts
	for _, u := range urls { // process the urls
		if safePool.processed[u] == 0 { // if we haven't processed the url yet
			wg.Add(1)           // add a task
			go func(u string) { // inline function takes an argument so that all goroutines don't share the same u
				defer wg.Done()                      // when the crawl finishes, we're done
				Crawl(u, depth-1, fetcher, safePool) // call Crawl recursively with one less level of depth
			}(u)
		}
	}
	safePool.mu.Unlock() // unlock before waiting, so other goroutines can make use of the pool
	wg.Wait()            // wait for other routines to finish
}

func main() {
	safePool := SafePool{processed: make(map[string]int)} // create a safepool for all goroutines to share
	Crawl("https://golang.org/", 4, fetcher, &safePool)   // start the crawl
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
