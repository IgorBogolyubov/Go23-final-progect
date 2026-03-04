package handlers

import (
	"Go23-final-progect/pkg/db"
	"Go23-final-progect/pkg/nextdate"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var task db.Task
	var resp db.Resp

	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		//Код с ошибкой выставляется внутри функции writeJson(w, resp)
		resp.Error = "Ошибка чтения формы запроса"
		writeJson(w, resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {

		resp.Error = "Ошибка десиреализации"
		writeJson(w, resp)
		return
	}

	if task.Title == "" {

		resp.Error = "Не задан заголовок"
		writeJson(w, resp)
		return

	}

	err = checkDate(&task)

	if err != nil {

		resp.Error = "Ошибка проверки даты"
		writeJson(w, resp)
		return
	}

	intId, err := db.AddTask(&task)

	if err != nil {

		resp.Error = "Ошибка конвнртации"
		writeJson(w, resp)
		return
	}
	task.ID = strconv.Itoa(int(intId))

	resp.ID = task.ID
	writeJson(w, resp)

}

func checkDate(task *db.Task) error {

	now := time.Now()
	var next string
	nowStr := now.Format("20060102")
	now, err := time.Parse("20060102", nowStr)
	if err != nil {

		return err
	}

	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {

		return err
	}

	if len(task.Repeat) != 0 {
		mass_interval := strings.Split(strings.TrimSpace(task.Repeat), " ")
		if mass_interval[0] != "d" && mass_interval[0] != "y" {

			return errors.New("не верный формат повторений")
		}
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {

			return errors.New("ошибка получения даты")
		}
	}

	if nextdate.AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			// next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var task db.Task
	var repeat []string
	var resp db.Resp
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		resp.Error = "Ошибка чтения формы запроса"
		writeJson(w, resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {

		resp.Error = "Ошибка десиреализации"
		writeJson(w, resp)
		return
	}

	if task.Title == "" {

		resp.Error = "Не задан заголовок"
		writeJson(w, resp)
		return

	}

	if len(task.Repeat) != 0 {
		repeat = strings.Split(strings.TrimSpace(task.Repeat), " ")

		if repeat[0] != "d" && repeat[0] != "y" {
			resp.Error = "не верный формат повторений"
			writeJson(w, resp)
			return

		}
	}

	err = checkDate(&task)

	if err != nil {

		resp.Error = "Ошибка проверки даты"
		writeJson(w, resp)
		return
	}

	/*
		Обработчик должен возвращать JSON с полем id или error. В случае успеха возвращается идентификатор созданной записи, а в случае ошибки — текст ошибки.

		{"id":"186"}

		{"error":"Не указан заголовок задачи"}

		Вот в каких случаях должны возвращаться ошибки:

		    ошибка десериализации JSON;
		    не указан заголовок задачи;
		    дата представлена в формате, отличном от 20060102;
		    правило повторения указано в неправильном формате.
	*/

	err = db.UpdateTask(&task)

	if err != nil {

		resp.Error = "Ошибка обновления"
		writeJson(w, resp)
		return
	}

	writeJson(w, resp)

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	type Resp struct {
		Error string `json:"error"`
	}
	resp1 := &Resp{}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")

	resp, err := db.GetTask(id)

	if err != nil {
		resp1.Error = "Задача не найдена"
		writeJson(w, resp1)
		return
	}
	writeJson(w, resp)

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	type TasksResp struct {
		Tasks string `json:"tasks,omitempty"`
	}
	type Resp struct {
		Error string `json:"error"`
	}
	resp := &Resp{}

	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)

	if err != nil {
		resp.Error = "Задача не найдена"
		writeJson(w, resp)
		return
	}

	writeJson(w, TasksResp{})
}
