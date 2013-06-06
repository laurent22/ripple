package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func HttpGet(url string) {
	log.Printf("GET %s\n", url)
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("    => " + string(body))
}

func HttpPost(url string, jsonString string) {
	log.Printf("POST %s\n", url)
	d := strings.NewReader(jsonString)
	response, _ := http.Post(url, "application/json", d)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("    => " + string(body))
}

func HttpPut(url string, jsonString string) {
	log.Printf("PUT %s\n", url)
	d := strings.NewReader(jsonString)
	request, _ := http.NewRequest("PUT", url, d)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("    => " + string(body))
}

func main() {
	baseUrl := "http://127.0.0.1:8080/"

	HttpGet(baseUrl + "users")
	HttpGet(baseUrl + "users/2")

	HttpPost(baseUrl+"users", "{\"Id\":0, \"Name\":\"New user\"}")
	HttpGet(baseUrl + "users")

	HttpPut(baseUrl+"users/5", "{\"Id\":0, \"Name\":\"Different name\"}")
	HttpGet(baseUrl + "users/5")

	HttpGet(baseUrl + "users/1/friends")
	HttpPost(baseUrl+"users/1/friends", "4")
	HttpGet(baseUrl + "users/1/friends")
}
