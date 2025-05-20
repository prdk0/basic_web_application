package main

import (
	"fmt"
	"net/http"
)

const PORT = ":8080"

func main() {
	// fmt.Println("Hello World")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintln(w, "Hello World!")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Number of bytes written %d\n", n)
	})
	fmt.Printf("Sever listening to the port %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
	}
}
