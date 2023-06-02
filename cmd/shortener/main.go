package main

import (
	"github.com/gorilla/mux"
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

func shortenURL(w http.ResponseWriter, r *http.Request) {

	url, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
	}

	urlID := Base62Encode(rand.Uint64())
	urlList[urlID] = string(url)
	response := "http://" + r.Host + "/" + urlID

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte(response))
}

func getFullURL(w http.ResponseWriter, r *http.Request) {

	shortURL := mux.Vars(r)["id"]

	val, ok := urlList[shortURL]
	if !ok {
		w.WriteHeader(400)
	}

	w.Header().Set("Location", val)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func main() {

	mux := mux.NewRouter()

	mux.HandleFunc(`/`, shortenURL).Methods("POST")
	mux.HandleFunc(`/{id}`, getFullURL).Methods("GET")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
