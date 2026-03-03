package api

import (
	"Go23-final-progect/handlers"
	"net/http"
)

const webDir = "./web"

func Init() {

	http.NewServeMux()

	http.HandleFunc("/api/task/done", handlers.DoneHandler)
	http.HandleFunc("/api/tasks", handlers.TasksHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler)
	http.HandleFunc("/api/nextdate", handlers.NextDayHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
