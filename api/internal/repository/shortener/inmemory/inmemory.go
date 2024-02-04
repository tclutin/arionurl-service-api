package inmemory

import (
	"errors"
	"github.com/tclutin/ArionURL/internal/service/shortener"
	"sync"
)

type shortenerMemoryRepo struct {
	sync.Mutex
	db map[string]*shortener.URL
}

func NewShortenerMemoryRepo() *shortenerMemoryRepo {
	return &shortenerMemoryRepo{
		db: make(map[string]*shortener.URL),
	}
}

func (d *shortenerMemoryRepo) Set(key string, model *shortener.URL) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.db[key] = model
}

func (d *shortenerMemoryRepo) Get(key string) (*shortener.URL, error) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	value, ok := d.db[key]
	if !ok {
		return nil, errors.New("error getting the value")
	}

	return value, nil
}

func (d *shortenerMemoryRepo) Delete(key string) error {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	if _, ok := d.db[key]; !ok {
		return errors.New("error getting the value")
	}

	delete(d.db, key)
	return nil
}
