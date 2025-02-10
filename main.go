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

	// TODO migrate to POST request when receiving the image from client
	http.HandleFunc("/checkbox", controller.CheckboxHandler)

	// Render api documentation
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
