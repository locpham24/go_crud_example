package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"main/route"
)

func main() {
	log.Println("Server started on: localhost:8080")
	router := route.Init()
	router.Logger.Fatal(router.Start(":8080"))
}
