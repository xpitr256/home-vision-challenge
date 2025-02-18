package main

import (
	_ "image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/xpitr256/home-vision-challenge/controller"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT environment variable not set, defaulting to %s", port)
	} else {
		log.Printf("Server will run on port %s", port)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/checkbox", controller.CheckboxHandler)

	// Render api documentation
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	// Allow clients to access response image with colored checkboxes
	fsResponse := http.FileServer(http.Dir("./response"))
	mux.Handle("/response/", http.StripPrefix("/response/", fsResponse))

	log.Printf("All routes configured")

	handler := enableCORS(mux)

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
