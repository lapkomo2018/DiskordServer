package main

import (
	"github.com/lapkomo2018/DiskordServer/internal/pkg/app"
	"log"
)

func main() {
	a, err := app.New(3000)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(a.Run())
}
