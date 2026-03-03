package main

import (
	"Go23-final-progect/pkg/api"
	"Go23-final-progect/pkg/db"
	"fmt"
)

const webDir = "./web"
const DbFile = "scheduler.db"

func main() {

	err := db.Init(DbFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	api.Init()

}
