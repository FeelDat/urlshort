package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345678"

var urlList = map[string]string{}

// Base62Encode Url shortening realisation
func Base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}
	return encodedBuilder.String()
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {

	url, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
	}

	urlId := Base62Encode(rand.Uint64())
	urlList[urlId] = string(url)
	response := "http://" + r.Host + "/" + urlId

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(response))
}

func getFullUrl(w http.ResponseWriter, r *http.Request) {

	shortUrl := r.Header.Get("id")

	val, ok := urlList[shortUrl[1:]]
	if !ok {
		w.WriteHeader(400)
	}

	w.Header().Set("Location", val)
	w.WriteHeader(307)

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortenUrl)
	mux.HandleFunc(`/{id}`, getFullUrl)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
