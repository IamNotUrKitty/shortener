package internal

import (
	"github.com/iamnoturkkitty/shortener/internal/config"
)

func Run() error {
	config := config.GetConfig()

	server, err := NewServer(config)
	if err != nil {
		return err
	}

	err = server.Start(config.Address)

	if err != nil {
		return err
	}

	return nil
}
