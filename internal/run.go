package internal

import (
	"github.com/iamnoturkkitty/shortener/internal/config"
)

func Run() error {
	config := config.GetConfig()

	server := NewServer(config)

	err := server.Start(config.Address)

	if err != nil {
		return err
	}

	return nil

}
