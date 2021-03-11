package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

func main() {
	//url := "192.168.1."
	argsWithProg := os.Args
	url := argsWithProg[1]
	var protocol = []string{"http", "https"}
	var port = []string{"80", "8080", "8081"}
	var wg sync.WaitGroup

	for i := 0; i < 256; i++ {
		host := i
		for _, x := range protocol {
			for _, b := range port {
				wg.Add(1)

				urls := fmt.Sprintf("%s://%s%d:%s", x, url, host, b)
				//fmt.Println(urls)
				go gettitle(urls, &wg)

			}
		}
	}
	wg.Wait()

}

func gettitle(urls string, wg *sync.WaitGroup) {
	defer func() {

	}()

	timeout := time.Duration(3 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(urls)

	if err != nil {

	} else {
		defer resp.Body.Close()
		bodyContent, _ := ioutil.ReadAll(resp.Body)

		var gettitle = regexp.MustCompile(`<title>(.*)</title>`)
		titleMatch := gettitle.FindStringSubmatch(string(bodyContent))
		title := "unknow"
		if len(titleMatch) != 0 {
			title = titleMatch[1]
		}
		fmt.Println(urls, resp.StatusCode, title)
	}

	wg.Done()

}
