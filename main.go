package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/xpitr256/home-vision-challenge/controller"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/checkbox", controller.CheckboxHandler)

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
