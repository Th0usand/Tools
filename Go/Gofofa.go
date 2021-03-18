package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Address struct {
	Mode    string        `json: mode`
	Error   bool          `json: error`
	Query   string        `json: query`
	Page    int           `json: page`
	Size    int           `json: size`
	Results []interface{} `json: results`
}

func main() {
	geturl()
}

func geturl() {
	//------------------------------------------------Modify 3 places--------------------------------------
	str_url1 := `app="tomcat"`
	email := ""
	key := ""
	
	
  
	strbytes := []byte(str_url1)
	str_url := base64.StdEncoding.EncodeToString(strbytes)
	url := "https://fofa.so/api/v1/search/all?email=" + email + "&key=" + key + "&qbase64=" + str_url + "&size=9999"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("get request failed, err:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	var addressObj Address
	err = json.Unmarshal(bodyContent, &addressObj)
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range addressObj.Results {
		a := value.([]interface{})[0].(string)

		if strings.Contains(a, "https") {
			b := a
			fmt.Println(b)
			write(b, str_url1)

		} else {
			b := "http://" + a
			fmt.Println(b)
			write(b, str_url1)
		}

	}
	return
}
func write(a string, str_url1 string) {
	filename := str_url1 + ".txt"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println(err)

	} else {

		_, err = f.Write([]byte(a + "\n"))
		if err != nil {
			fmt.Println(err)

		}
	}
}
