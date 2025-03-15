package main

import (
	"fmt"
	"github.com/karkiayush/Stock-CRUD-API/router"
	"log"
	"net/http"
)

func main() {
	const PORT = 8080
	r := router.Router()

	fmt.Println("Starting Server on port: ", PORT)
	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", PORT),
		r,
	); err != nil {
		log.Fatal(err)
	}
}
