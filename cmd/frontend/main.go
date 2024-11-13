package main

import "net/http"

func main() {
	if err := http.ListenAndServe(""); err != nil {
		panic(err)
	}
}