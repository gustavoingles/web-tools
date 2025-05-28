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

func main() {
	var wg sync.WaitGroup
	urls := getUrls()

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

func urlParser(url *string) string {
	removedScheme := strings.TrimPrefix(*url, "https://")
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
