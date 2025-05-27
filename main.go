package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	urls := []string{"https://www.google.com", "https://pkg.go.dev/std", "https://www.github.com"}

	wg.Add(len(urls))
	for _, v := range urls {
		go fetchUrl(v, &wg)
	}
	wg.Wait()
}

func fetchUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("%s is a valid URL\n", url)
	} else {
		log.Printf("%s returned status code %d\n", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing response's body: %v\n", err)
	}

	fmt.Println(string(body))
}
