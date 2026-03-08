package handlers

import (
	"Go23-final-progect/pkg/db"
	"Go23-final-progect/pkg/nextdate"
	"net/http"
	"strconv"
	"time"
)

func NextDayHandler(w http.ResponseWriter, req *http.Request) {

	var now time.Time
	var err error

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(req.FormValue("now")) == 0 {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", req.FormValue("now"))

		if err != nil {
			http.Error(w, "Error parse", http.StatusBadRequest)
			return
		}
	}

	date := req.FormValue("date")

	repeat := req.FormValue("repeat")

	answer, err := nextdate.NextDate(now, date, repeat)

	if err != nil {
		http.Error(w, "Error date", http.StatusBadRequest)
		return
	}

	w.Write([]byte(answer))
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodGet:
		GetTaskHandler(w, r)

	case http.MethodPost:
		AddTaskHandler(w, r)

	case http.MethodPut:
		UpdateTaskHandler(w, r)

	case http.MethodDelete:
		DeleteTaskHandler(w, r)

	default:
		// Возвращаем статус 405, если метод не поддерживается
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	type TasksResp struct {
		Tasks []*db.Task `json:"tasks"`
	}
	type errResp struct {
		Error string `json:"error"`
	}
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		writeJson(w, errResp{Error: "Ошибка запроса"})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func DoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	type TasksResp struct {
		Tasks string `json:"tasks,omitempty"`
	}

	type RespE struct {
		Error string `json:"error"`
	}
	respError := &RespE{}

	id := r.URL.Query().Get("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadRequest)
		respError.Error = "Не верно задан ID"
		writeJson(w, respError)
		return
	}
	resp, err := db.GetTask(id)

	if err != nil {
		http.Error(w, "Task not faund", http.StatusBadRequest)
		respError.Error = "Задача не найдена"
		writeJson(w, respError)
		return
	}

	if len(resp.Repeat) != 0 {

		now, err := time.Parse("20060102", resp.Date)
		if err != nil {
			http.Error(w, "error parse", http.StatusBadRequest)
			respError.Error = "ошибка даты"
			writeJson(w, respError)
			return
		}

		answer, err := nextdate.NextDate(now, resp.Date, resp.Repeat)
		if err != nil {
			http.Error(w, "Bad Request data", http.StatusBadRequest)
			respError.Error = "ошибка получения даты"
			writeJson(w, respError)
			return
		}

		err = db.UpdateDate(answer, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			respError.Error = "Ошибка обновления даты"
			writeJson(w, respError)
			return
		}

	} else {
		err = db.DeleteTask(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			respError.Error = "Ошибка удаления записи"
			writeJson(w, respError)
			return
		}
	}

	writeJson(w, TasksResp{})
}
