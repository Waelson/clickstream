package main

import (
	"log"
	"net/http"
)

func main() {
	// Servir arquivos estaticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Iniciar o servidor
	port := ":8080"
	log.Printf("items-frontend running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
