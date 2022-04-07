// Запустить сервис.

package main

import (
	"log"
    "github.com/opnely/WB_T3/cmd/service"
)

func main() {
	router, err := NewRouter()
	if err != nil {
		log.Fatalf("NewRouter: %v\n", err)
	}
	router.Start()
}
