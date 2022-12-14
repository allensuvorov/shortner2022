package storage

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

type urlStorage struct {
	inMemory inMemory
}

type inMemory struct {
	URLHashMap     hashmap.URLHashMap
	ClientActivity hashmap.ClientActivity
	Deleted        hashmap.Deleted
}

func NewURLStorage() *urlStorage {
	fsp := config.UC.FSP
	um := make(hashmap.URLHashMap)
	ca := make(hashmap.ClientActivity)
	dd := make(hashmap.Deleted)
	im := inMemory{um, ca, dd}

	if fsp != "" {
		im = restore(fsp)
	}

	return &urlStorage{
		inMemory: im,
	}
}

func (us *urlStorage) Create(ue entity.URLEntity) error {
	log.Println("Storage/Create(): hello")

	us.inMemory.URLHashMap[ue.Hash] = ue.URL
	log.Println("Storage/Create(): added to map, updated map len is", len(us.inMemory.URLHashMap))
	log.Println("Storage/Create(): added to map, updated map is", us.inMemory.URLHashMap)

	_, ok := us.inMemory.ClientActivity[ue.ClientID]
	if !ok {
		us.inMemory.ClientActivity[ue.ClientID] = make(map[string]bool)
	}
	us.inMemory.ClientActivity[ue.ClientID][ue.Hash] = true

	fsp := config.UC.FSP

	if fsp != "" {
		write(ue, fsp)
	}
	log.Printf("Storage/Create(): created hash: %s, for URL: %s. File path %s:", ue.Hash, ue.URL, fsp)
	return nil
}

func (us *urlStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage/GetHashByURL(), looking for matching URL", u)
	for k, v := range us.inMemory.URLHashMap {
		if v == u {
			log.Println("Storage GetHashByURL, found record", k)
			return k, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *urlStorage) GetURLByHash(h string) (string, error) {
	log.Println("Storage/GetURLByHash(), looking in map len", len(us.inMemory.URLHashMap))
	log.Println("Storage/GetURLByHash(), looking for matching Hash", h)
	u, ok := us.inMemory.URLHashMap[h]
	if !ok {
		return "", errors.ErrNotFound
	}
	if us.inMemory.Deleted[h] {
		return "", errors.ErrRecordDeleted
	}
	return u, nil
}

func (us *urlStorage) GetClientUrls(id string) ([]entity.URLEntity, error) {
	log.Println("storage/GetClientUrls client id is:", id)
	ca, ok := us.inMemory.ClientActivity[id]
	if !ok {
		return nil, nil
	}
	log.Println("storage/GetClientUrls client ClientActivity is:", ca)
	dtoList := []entity.URLEntity{}

	for k := range ca {
		u, err := us.GetURLByHash(k)
		bu := config.UC.BU
		if err != nil {
			return nil, err
		}
		ue := entity.URLEntity{
			Hash: bu + "/" + k,
			URL:  u,
		}
		dtoList = append(dtoList, ue)
	}
	log.Println("storage/GetClientUrls dtoList is:", dtoList)

	return dtoList, nil
}

func (us *urlStorage) PingDB() bool {
	return true
}

func (us *urlStorage) BatchDelete(hashList *[]string, clientID string) error {
	log.Println("Storage/BatchDelete - Hello")

	for _, h := range *hashList {
		_, ok := us.inMemory.ClientActivity[clientID][h]
		if ok {
			us.inMemory.Deleted[h] = true
		}
	}
	log.Println("Storage/BatchDelete - inMemory.URLHashMap:", us.inMemory.URLHashMap)
	log.Println("Storage/BatchDelete - inMemory.ClientActivity:", us.inMemory.ClientActivity)
	log.Println("Storage/BatchDelete - inMemory.Deleted:", us.inMemory.Deleted)
	log.Println("Storage/BatchDelete - Bye")
	return nil
}
