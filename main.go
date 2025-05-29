package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yosssi/gohtml"
)

var rightCount atomic.Int64

var client = http.Client{
	Timeout: 30 * time.Second,
}

func main() {
	var wg sync.WaitGroup
	urls := make(chan string, 100)
	workerCount := runtime.NumCPU() * 2
	for range workerCount  {
		wg.Add(1)
		go fetchUrl(urls, &wg)
	}
	go getUrls(urls)
	wg.Wait()
	fmt.Printf("got %d requests right", rightCount.Load())
}

func fetchUrl(urls <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urls {
		func ()  {
			resp, err := client.Get(url)
			if err != nil {
				log.Printf("error fetching %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Printf("%s is a valid URL", url)
				rightCount.Add(1)
			} else {
				log.Printf("%s returned status code %d", url, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error parsing response's body: %v", err)
				return
			}

			// Check resp.Request.Url data about the whole URL and its pieces
			// fileName := resp.Request.URL.Hostname()
			fileName := resp.Request.URL.Hostname()

			err = writeToFile(fileName, body)
			if err != nil {
				log.Printf("error creating, opening, or writing into the file: %v", err)
				return
			}
		} ()


	}
}

func getUrls(urls chan<- string) {
	defer close(urls)
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			urls <- os.Args[i]
		}
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls <- line
		}
	}
}

func writeToFile(fileName string, fileContent []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	formattedFile := gohtml.FormatBytes(fileContent)
	_, err = file.Write(formattedFile)
	return err
}
