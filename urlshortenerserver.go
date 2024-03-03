package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var (
	shortpath   = "http://localhost:9000/"
	urlnotfound = "url not found"
	keyLength   = 7
	PORT        = "9009"
)

type Request struct {
	Url string
}

type Response struct {
	Shorturl string
}

type UrlStore struct {
	urls map[string]string
}

func (url *UrlStore) CreateUrl(urlToReplace string) (string, error) {
	key := GenerateKey()
	shorturl := shortpath + key
	url.AddUrlToStore(key, urlToReplace)
	return shorturl, nil
}

func (url *UrlStore) AddUrlToStore(key string, newUrl string) {
	url.urls[key] = newUrl
}

func (url *UrlStore) GetUrl(key string) string {
	if url, ok := url.urls[key]; ok {
		return url
	} else {
		return urlnotfound
	}
}

func (url *UrlStore) RemoveUrl(keyToRemove string) {
	delete(url.urls, keyToRemove)
}

func GenerateKey() string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	key := make([]string, keyLength)
	for i := 0; i < keyLength; i++ {
		randomNumber := rand.Intn(len(alphabet))
		v := alphabet[randomNumber]
		key = append(key, string(v))
	}
	str := strings.Join(key, "")

	return str
}

func (url *UrlStore) URLShorteningRequesthandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	var reqBody Request
	err := decoder.Decode(&reqBody)

	if err != nil {
		fmt.Println("something went wrong", err)
		defer io.WriteString(res, `{"status":"failed"}`)
	}

	newUrl, _ := url.CreateUrl(reqBody.Url)
	responseObj := &Response{Shorturl: newUrl}
	response, err := json.Marshal(responseObj)

	if err != nil {
		fmt.Println("something went wrong", err)
		defer io.WriteString(res, `{"status":"failed"}`)
	}

	defer io.WriteString(res, string(response))

}

func (url *UrlStore) Redirect(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Received redirect request")
	key := req.URL.Path[1:]
	fmt.Println("Key..", string(key))
	intitialUrl := url.GetUrl(string(key))
	http.Redirect(res, req, intitialUrl, http.StatusFound)
}

func main() {
	urlShortner := &UrlStore{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shortenurl", urlShortner.URLShorteningRequesthandler)
	http.HandleFunc("/", urlShortner.Redirect)
	fmt.Println("Server is up on port ", PORT)
	http.ListenAndServe(PORT, nil)
}
