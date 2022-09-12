package url

import (
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"

)

type URLStorage interface {

}

type URLService struct {
	urlStorage URLStorage
}

func NewURLService(us URLStorage) URLService {
	return NewURLService{
		urlStorage: us,
	}
}
// func CreateURL

// check if URL is valid
_, err = url.ParseRequestURI(u)

if err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

// get Hash if the longURL already exists in storage
h, ok := storage.GetHash(u)

// if longURL does not exist in storage
if !ok {

	// generate shortened URL
	h = Shorten(u)

	// add url to the storage
	storage.CreateHash(h, u)
}

shortURL := "http://localhost:8080/" + h