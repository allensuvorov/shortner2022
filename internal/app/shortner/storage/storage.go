package storage

import (
	"log"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

// Object with storage methods to work with DB
type URLStorage struct {
	InMemory InMemory
}

type InMemory struct {
	URLHashMap   hashmap.URLHashMap
	UserActivity hashmap.UserActivity
}

// NewURLStorage creates URLStorage object
func NewURLStorage() *URLStorage {
	// Restore data at start up
	fsp := config.UC.FSP
	um := make(hashmap.URLHashMap) // url map
	ua := make(hashmap.UserActivity)

	// TODO restore from file both um and ua

	// restore if path in config not empty
	if fsp != "" {
		um = restore(fsp) // get map
	}

	return &URLStorage{
		InMemory: InMemory{um, ua},
	}
}

// Create adds new URL record to storage
func (us *URLStorage) Create(h, u string) error {
	// Save to map
	us.InMemory.URLHashMap[h] = u
	log.Println("Storage/Create(): added to map, updated map len is", len(us.InMemory.URLHashMap))
	log.Println("Storage/Create(): added to map, updated map is", us.InMemory.URLHashMap)

	// get file storage path from config
	fsp := config.UC.FSP

	// Save to file, if there is path in config
	if fsp != "" {
		write(h, u, fsp)
	}
	log.Printf("Storage/Create(): created hash: %s, for URL: %s. File path %s:", h, u, fsp)
	return nil
}

func (us *URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage/GetHashByURL(), looking for matching URL", u)
	for k, v := range us.InMemory.URLHashMap {
		if v == u {
			log.Println("Storage GetHashByURL, found record", k)
			return k, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *URLStorage) GetURLByHash(h string) (string, error) {
	log.Println("Storage/GetURLByHash(), looking in map len", len(us.InMemory.URLHashMap))
	log.Println("Storage/GetURLByHash(), looking for matching Hash", h)
	u, ok := us.InMemory.URLHashMap[h]
	if !ok {
		return "", errors.ErrNotFound
	}
	return u, nil
}
