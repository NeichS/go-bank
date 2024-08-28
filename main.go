package main

import (
	"fmt"
	"log"
)

func main() {

	store, err := NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Conectado a la base de datos exitosamente")
	}

	api := NewAPIServer(":3000", store)
	api.Run()
}
