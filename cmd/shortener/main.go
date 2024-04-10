package main

import (
	"github.com/iamnoturkkitty/shortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
