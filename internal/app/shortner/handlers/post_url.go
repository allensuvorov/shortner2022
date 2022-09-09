package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/util"
)

// PostURL - takes full URL and returns short URL.
func PostURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if long URL is empty string
	if len(b) == 0 {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}
	// get u - string of b
	u := string(b)

	// check if URL is valid
	_, err = url.ParseRequestURI(u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get Hash and check if long url is already in the map
	h, err := storage.GetHash(u)

	// log body from request
	log.Println(u)

	// if Hash already exists
	if err == nil {

		// generate shortened URL
		h = util.Shorten(u)

		// add url to the map
		storage.CreateHash(h, u)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	shortURL := "http://localhost:8080/" + h

	// пишем тело ответа
	w.Write([]byte(shortURL))

}
