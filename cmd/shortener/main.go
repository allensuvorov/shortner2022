package main

import (
	"io"
	"log"
	"net/http"
)

var urls map[string]string = make(map[string]string)

// CreateShortURL — обработчик запроса.
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
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
		urls[bs] = bs
		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write(b)
	} else {
		w.Write([]byte("<h1>Hello, World</h1>"))
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", CreateShortURL)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
