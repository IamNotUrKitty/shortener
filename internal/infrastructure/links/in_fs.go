package links

import (
	"encoding/json"
	"io"
	"os"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type InFSRepo struct {
	memory  *InMemoryRepo
	file    *os.File
	encoder *json.Encoder
}

func readLink(decoder *json.Decoder) (*links.StoredLink, error) {
	l := &links.StoredLink{}
	if err := decoder.Decode(&l); err != nil {
		return nil, err
	}

	return l, nil
}

func loadLinksFromFS(mr *InMemoryRepo, file *os.File) error {
	decoder := json.NewDecoder(file)
	for {
		e, err := readLink(decoder)
		if err == io.EOF {
			break
		}

		l, err := links.NewLink(e.ID, e.URL, e.Hash)
		if err != nil {
			return err
		}

		if err := mr.SaveLink(*l); err != nil {
			return err
		}
	}

	return nil
}

func NewInFSRepo(fileName string) (*InFSRepo, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	memRep := NewInMemoryRepo()

	if err := loadLinksFromFS(memRep, file); err != nil {
		return nil, err
	}

	return &InFSRepo{
		file:    file,
		memory:  memRep,
		encoder: json.NewEncoder(file),
	}, nil
}

func (r *InFSRepo) WriteLink(l links.Link) error {
	return r.encoder.Encode(links.StoredLink{ID: l.ID(), URL: l.URL(), Hash: l.Hash()})
}

func (r *InFSRepo) SaveLinkBatch(ls []links.Link) error {
	for _, l := range ls {
		r.SaveLink(l)
	}

	return nil
}

func (r *InFSRepo) SaveLink(l links.Link) error {
	if err := r.memory.SaveLink(l); err != nil {
		return err
	}

	if err := r.WriteLink(l); err != nil {
		return err
	}

	return nil
}

func (r *InFSRepo) GetLink(hash string) (*links.Link, error) {
	return r.memory.GetLink(hash)
}

func (r *InFSRepo) Test() error {
	return nil
}
