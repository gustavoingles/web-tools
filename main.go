package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yosssi/gohtml"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	var wg sync.WaitGroup
	urls := getUrls()

	for _, v := range urls {
		wg.Add(1)
		go fetchUrl(v, &wg)
	}
	wg.Wait()
}

func fetchUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("%s is a valid URL", url)
	} else {
		log.Printf("%s returned status code %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error parsing response's body: %v", err)
	}

	// Check resp.Request.Url data about the whole URL and its pieces
	// fileName := resp.Request.URL.Hostname()
	fileName := urlParser(url)

	err = writeToFile(fileName, body)
	if err != nil {
		log.Printf("error creating, opening, or writing into the file: %v", err)
	}
}

func getUrls() []string {
	if len(os.Args) > 1 {
		urls := os.Args[1:]
		return urls
	} 
	urls := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls
}

func urlParser(url string) string {
	removedScheme := strings.TrimPrefix(url, "https://")
	removedTLD := strings.TrimSuffix(removedScheme, ".com")
	domainHTML := removedTLD + ".html"

	return domainHTML
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
