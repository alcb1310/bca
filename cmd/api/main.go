package main

import (
	"bca-go-final/internal/database"
	"bca-go-final/internal/server"
)

func main() {
	db := database.New()
	server, _ := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
