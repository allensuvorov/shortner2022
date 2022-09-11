package storage

import (
	"log"
	"yandex/projects/urlshortner/internal/app/shortner/domain/entity"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	//"github.com/allensuvorov/urlshortner/internal/app/domain/entity"
)

// map to store short urls and full urls
// var Urls map[string]string = make(map[string]string)

// DBStorage interface (to be adopted for mock-testing):
type DBStorage interface {

	// check if hash exists
	HashExists(u string) bool

	// return longURL for the matching hash
	GetURL(u string) string

	// return hash for a matching longURL, check if longURL exists
	GetHash(u string) (string, error)

	// add new record - pair shortURL: longURL
	CreateHash(h, string, u string) error
}

// Object with storage methods to work with DB
type URLStorage struct {
	inMemory []entity.URLEntity
}

// NewURLStorage creates URLStorage object
func NewURLStorage() URLStorage {
	return URLStorage{
		inMemory: make([]entity.URLEntity, 0),
	}
}

// Create adds new url record to storage
func (us URLStorage) Create(url entity.URLEntity) error {
	us.inMemory = append(us.inMemory, url)

}

// add new record - pair shortURL: longURL
func CreateHash(h string, u string) error {
	Urls[h] = u
	return nil
}

// Get

// HashExists checks if hash exists.
func HashExists(h string) bool {
	if _, ok := Urls[h]; ok {
		return true
	}
	return false
}

// GetHash returns hash for the original longURL.
func GetHash(u string) (string, bool) {
	for k, v := range Urls {
		if v == u {
			log.Println("URL already exists")
			return k, true
		}
	}
	return "", false
}

// GetURL returns longURL for the matching hash.
func GetURL(h string) (string, bool) {
	if u, ok := Urls[h]; ok {
		return u, true
	}
	return "", false
}
