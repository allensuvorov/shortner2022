package url

import (
	"io"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
)

type URLService interface {
	// Create takes URL and returns hash.
	Create(u string) (string, error)
	// Get takes hash and returns URL.
	Get(h string) (entity.URLEntity, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(us URLService) URLHandler {
	return URLHandler{
		urlService: us,
	}
}

// Create passes URL to service and returns response with Hash.
func (uh URLHandler) Create(w http.ResponseWriter, r *http.Request) {

	// читаем Body
	b, err := io.ReadAll(r.Body)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert to string
	u := string(b)

	// check if long URL is empty string
	if len(b) == 0 {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}

	// log body from request
	log.Println("URL in the POST request is", u)

	shortURL, err := uh.urlService.Create(u)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	w.Write([]byte(shortURL))
}

// Get passes Hash to service and returns response with URL.
func (uh URLHandler) Get(w http.ResponseWriter, r *http.Request) {

}