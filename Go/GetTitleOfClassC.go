package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

func httpGet(url string) (err error) {
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)

	var gettitle = regexp.MustCompile(`<title>(.*)</title>`)
	titleMatch := gettitle.FindStringSubmatch(string(bodyContent))
	title := "unknow"
	if len(titleMatch) != 0 {
		title = titleMatch[1]
	}
	fmt.Println(url, resp.StatusCode, title)
	return
}

func RequestPort(protocol string, host string, port string) {
	wg := sync.WaitGroup{}
	for i := 0; i < 256; i++ {
		urls := fmt.Sprintf("%s://%s%d:%s", protocol, host, i, port)
		wg.Add(1)
		go func(url string) {

			httpGet(url)
			wg.Done()
		}(urls)
	}
	wg.Wait()
	return
}

func main() {
	var url = "181.198.240."
	ports := []string{"80", "443", "8080"}
	wg := sync.WaitGroup{}
	for _, i := range ports {
		wg.Add(1)
		go func(port string) {
			RequestPort("http", url, port)
			wg.Done()
		}(i)
	}
	wg.Wait()

}
