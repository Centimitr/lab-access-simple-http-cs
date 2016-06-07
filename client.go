package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Message  string `json:"message"`
	NextId   int    `json:"id,omitempty"`
	NextCode string `json:"code,omitempty"`
}

func getContent(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func getMessage(url string) Message {
	content := getContent(url)
	var m Message
	json.Unmarshal([]byte(content), &m)
	return m
}

func do() {
	baseurl := "http://10.10.65.110:3000"
	str := getContent(baseurl)
	nid, _ := strconv.Atoi(strings.Split(str, " ")[0])
	m := Message{
		NextId:   nid,
		NextCode: strings.Split(str, " ")[1],
	}
	for m.Message != "Congratulations!" {
		fmt.Print("ASK ", strconv.Itoa(m.NextId), ": ")
		v := neturl.Values{}
		v.Set("id", strconv.Itoa(m.NextId))
		v.Set("code", m.NextCode)
		url := baseurl + "/ask/?" + v.Encode()
		m = getMessage(url)
		fmt.Println(m.Message)
	}
}

func main() {
	for i := 0; i < 10; i++ {
		go do()
	}
	time.Sleep(time.Minute)
}
