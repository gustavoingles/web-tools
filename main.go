package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	urls := []string{"https://google.com", "https://theodinproject.com", "https://github.com"}

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

	fileName := urlParser(&url) 

	err = writeToFile(fileName, body)
	if err != nil {
		log.Printf("Error creating, opening, or writing into the file: %v\n", err)
	}
}

func urlParser(url *string) string {
	removedScheme := strings.TrimPrefix(*url, "https://")
	removedTLD := strings.TrimSuffix(removedScheme, ".com")
	domainHTML := removedTLD + ".html"

	return domainHTML
}

func writeToFile(fileName string, fileContent []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(fileContent)
	return err
}
