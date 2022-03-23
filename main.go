// Запустить сервис

package main

import (
	"log"
)

func main() {
	router, err := NewRouter()
	if err != nil {
		log.Fatalf("NewRouter: %v\n", err)
	}
	router.MakeRoutes()
	log.Println("Запуск сервера")
	router.Start()
}
