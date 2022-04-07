// Запустить сервис.

package main

import (
	"log"
    "github.com/Opnely/WB_T3/cmd/service"
)

func main() {
	router, err := service.NewRouter()
	if err != nil {
		log.Fatalf("NewRouter: %v\n", err)
	}
	router.Start()
}
