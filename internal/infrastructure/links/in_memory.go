package links

import (
	"sync"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type InMemoryRepo struct {
	links map[string]string
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		links: make(map[string]string),
	}
}

func (r *InMemoryRepo) SaveLink(l links.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.links[l.Hash()] = l.Url()

	return nil
}

func (r *InMemoryRepo) GetLink(hash string) (*links.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, ok := r.links[hash]

	if !ok {
		return nil, links.LinkNotFound
	}
	link, err := links.NewLink(l)
	if err != nil {
		return nil, err
	}

	return link, nil
}
