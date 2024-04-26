package links

import (
	"encoding/json"
	"io"
	"os"

	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

type InFSRepo struct {
	memory  *InMemoryRepo
	file    *os.File
	encoder *json.Encoder
}

func InitFSRepo(fileName string) (handlers.Repository, error) {
	if fileName == "" {
		return NewInMemoryRepo(), nil
	}
	return NewInFSRepo(fileName)
}

func readLink(decoder *json.Decoder) (*links.StoredLink, error) {
	l := &links.StoredLink{}
	if err := decoder.Decode(&l); err != nil {
		return nil, err
	}

	return l, nil
}

func NewInFSRepo(fileName string) (*InFSRepo, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	memRep := NewInMemoryRepo()
	decoder := json.NewDecoder(file)

	for {
		e, err := readLink(decoder)
		if err == io.EOF {
			break
		}

		l, err := links.NewLink(e.ID, e.URL, e.Hash)
		if err != nil {
			return nil, err
		}

		memRep.SaveLink(*l)
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

func (r *InFSRepo) SaveLink(l links.Link) error {
	r.memory.SaveLink(l)

	r.WriteLink(l)

	return nil
}

func (r *InFSRepo) GetLink(hash string) (*links.Link, error) {

	return r.memory.GetLink(hash)
}