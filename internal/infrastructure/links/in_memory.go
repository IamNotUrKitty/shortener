package links

import (
	"context"
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

func (r *InMemoryRepo) SaveLink(ctx context.Context, l links.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.links[l.Hash()] = l

	return nil
}

func (r *InMemoryRepo) SaveLinkBatch(ctx context.Context, ls []links.Link) error {
	for _, l := range ls {
		r.SaveLink(ctx, l)
	}

	return nil
}

func (r *InMemoryRepo) GetLink(ctx context.Context, hash string) (*links.Link, error) {
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

func (r *InMemoryRepo) Test() error {
	return nil
}
