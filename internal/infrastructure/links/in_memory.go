package links

import (
	"fmt"
	"sync"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type InMemoryRepo struct {
	links map[string]links.Link
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		links: make(map[string]links.Link),
	}
}

func (r *InMemoryRepo) SaveLink(l links.Link) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	// r.links[l.Hash()] = l
	fmt.Println("123")

	return nil
}

func (r *InMemoryRepo) GetLink(hash string) (*links.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, ok := r.links[hash]

	if !ok {
		return nil, links.ErrLinkNotFound
	}
	link, err := links.NewLink(l.ID(), l.URL(), l.Hash())
	if err != nil {
		return nil, err
	}

	return link, nil
}
