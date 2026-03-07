package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Resp struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error add task: %w", err)
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	var res []*Task

	query := "SELECT * FROM scheduler ORDER BY Date LIMIT :limit"
	rows, err := db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ress := Task{}
		err = rows.Scan(&ress.ID, &ress.Date, &ress.Title, &ress.Comment, &ress.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error in scan: %w", err)
		}
		res = append(res, &ress)
	}
	err = rows.Err()

	if err != nil {

		return nil, fmt.Errorf("error iteration: %w", err)
	}

	if len(res) == 0 {
		res = []*Task{}
	}

	return res, nil

}

func UpdateTask(task *Task) error {

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("error update task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("data is not update: %w", err)
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func GetTask(id string) (*Task, error) {

	id_Int, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf(`incorrect id`)

	}
	query := "SELECT * FROM scheduler WHERE id = :id"
	row := db.QueryRow(query, sql.Named("id", id_Int))
	ress := Task{}

	err = row.Scan(&ress.ID, &ress.Date, &ress.Title, &ress.Comment, &ress.Repeat)
	if err != nil {
		return nil, fmt.Errorf(`incorrect task`)

	}

	return &ress, nil
}

func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("not ID")
	}

	id1, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("error ID: %w", err)
	}

	query := `DELETE FROM scheduler WHERE id = :id`

	_, err = db.Exec(query, sql.Named("id", id1))
	if err != nil {
		return fmt.Errorf("error delete task: %w", err)
	}

	return nil
}

func UpdateDate(next string, id string) error {
	if id == "" {
		return fmt.Errorf("not ID")
	}

	id_Int, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", next), sql.Named("id", id_Int))
	if err != nil {
		return fmt.Errorf("data is not update: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("data is not update: %w", err)
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id query`)
	}
	return nil
}
