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
	serviceUrl string
	wg         *sync.WaitGroup
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	numArgs := len(os.Args)
	if numArgs > 1 {
		f := &fetcher{
			serviceUrl: "https://www.googleapis.com/urlshortener/v1/url",
			wg:         new(sync.WaitGroup),
		}
		f.wg.Add(numArgs - 1)
		defer f.wg.Wait()
		for _, url := range os.Args[1:] {
			go f.fetch(url)
		}
	}
}

func (f *fetcher) fetch(url string) {
	client := new(http.Client)
	defer f.wg.Done()
	data := make(map[string]string)
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
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Printf("Could not unmarshal data with response: %s", string(response))
		return
	}
	fmt.Printf("%s ==> %s\n", url, data["id"])
}
