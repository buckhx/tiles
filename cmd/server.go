package main

import "github.com/buckhx/tiles"

func main() {
	server := tiles.NewServer(":8080")
	server.ListenAndServe()
}
