package main

import (
	"github.com/iamnoturkkitty/shortener/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		panic(err)
	}
}
