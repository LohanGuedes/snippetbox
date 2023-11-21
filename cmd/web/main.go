package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse() // Automatically quits if there's an error.

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet/view", SnippetView)
	mux.HandleFunc("/snippet/create", SnippetCreate)

	log.Println("Starting server on port", *addr)

	err := http.ListenAndServe("127.0.0.1"+*addr, mux) // The port is to be declared like this in MacOS in order to remove that allow popup
	log.Fatal(err)
}
