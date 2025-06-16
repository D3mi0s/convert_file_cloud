package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	log.Println("Frontend server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
