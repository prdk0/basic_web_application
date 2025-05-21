package main

import (
	"bwa/pkg/handlers"
	"fmt"
	"net/http"
)

const PORT = ":8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	fmt.Printf("Sever listening to the port %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
	}
}
