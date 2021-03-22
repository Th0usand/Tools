package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ResponseStruct struct {
	ResponseBody []byte
	Status       int
	Headers      http.Header
}

var wg sync.WaitGroup

func main() {
	url := "http://192.168.35."
	//argsWithProg := os.Args
	//url := argsWithProg[1]
	for i := 0; i < 256; i++ {
		host := i
		wg.Add(1)
		url1 := fmt.Sprintf("%s%d", url, host)
		go Getdav(url1)

	}
	wg.Wait()

}

func Getdav(url string) {
	defer wg.Done()
	resp, err := Request(url, "OPTIONS", time.Second*3)
	if err != nil {

	}

	values := fmt.Sprintf("%v", resp.Headers)
	if strings.Contains(values, "DAV") {
		fmt.Println("OpenWebdavIp:", url)
	}
}

func Request(url string, method string, timeout time.Duration) (respObj ResponseStruct, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ResponseStruct{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseStruct{}, err
	}
	defer resp.Body.Close()

	respObj.ResponseBody, _ = ioutil.ReadAll(resp.Body)
	respObj.Status = resp.StatusCode
	respObj.Headers = resp.Header
	return respObj, nil

}
