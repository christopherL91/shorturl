package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

type fetcher struct {
	// shorturl service url
	serviceUrl string
	// waitgroup for all the requests
	wg *sync.WaitGroup
}

func init() {
	// Run with maximum number of cores available
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// number of urls
	numUrls := len(os.Args) - 1
	if numUrls > 0 {
		f := &fetcher{
			serviceUrl: "https://www.googleapis.com/urlshortener/v1/url",
			wg:         new(sync.WaitGroup),
		}
		// Set number of goroutines to wait for
		f.wg.Add(numUrls)
		// wait for all the goroutines
		defer f.wg.Wait()
		for _, url := range os.Args[1:] {
			// launch new goroutine and fetch short url
			go f.fetch(url)
		}
	}
}

func (f *fetcher) fetch(url string) {
	client := new(http.Client)
	defer f.wg.Done()
	// response from server
	data := make(map[string]string)
	// prepare body for request
	body := bytes.NewBufferString(fmt.Sprintf(`{"longUrl":"%s"}`, url))
	res, err := client.Post(f.serviceUrl, "application/json", body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Response error with data: %s. Error:%s", url, err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Could not read response with data: %s", url)
		return
	}
	// put the response in data
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Printf("Could not unmarshal data with response: %s", string(response))
		return
	}
	fmt.Printf("%s ==> %s\n", url, data["id"])
}
