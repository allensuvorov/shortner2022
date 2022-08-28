package main

import (
	"io"
	"log"
	"net/http"
	"path"
)

// map to store short urls and full urls
var urls map[string]string = make(map[string]string)

// mock shortner
func Shorten(s string) string {
	return "shortURL:" + s[0:5]
}

// CreateShortURL — обработчик запроса.
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if r.Method == "GET" {

			base := path.Base(r.URL.Path)
			//http.Error(w, "404 not found.", http.StatusNotFound)

			// устанавливаем статус-код 307
			w.WriteHeader(http.StatusTemporaryRedirect)

			// set header location
			// ...

			w.Write([]byte(urls[base]))

			return
		}
	}
	if r.Method == "POST" {
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		bs := string(b) // body string

		// add url to the map
		urls[Shorten(bs)] = bs
		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte(Shorten(bs)))
	}

}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", CreateShortURL)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
