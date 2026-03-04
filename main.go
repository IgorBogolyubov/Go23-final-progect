package main

import (
	"Go23-final-progect/pkg/api"
	"Go23-final-progect/pkg/db"
	"log"
)

const webDir = "./web"
const DbFile = "scheduler.db"

func main() {

	err := db.Init(DbFile)
	if err != nil {
		log.Fatalf("error in init DB: %v", err)
	}

	api.Init()

}
