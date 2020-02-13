package model

import (
	"github.com/labstack/echo"
	"log"
	"main/db"
	"strconv"
)

type Task struct {
	Id     int
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
	Done   int    `json:"done"`
}

type Tasks []Task

func TaskAll(userId int) Tasks {
	task := Task{}
	res := Tasks{}

	db := db.GetInstance()
	selDb, err := db.Query("SELECT * FROM task WHERE user_id='" + strconv.Itoa(userId) + "' ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	for selDb.Next() {
		var id, userId, done int
		var title string
		err = selDb.Scan(&id, &userId, &title, &done)
		if err != nil {
			panic(err.Error())
		}
		task.Id = id
		task.UserId = userId
		task.Title = title
		task.Done = done

		res = append(res, task)
	}

	return res
}

func TaskOne(taskId int) Task {
	db := db.GetInstance()

	selDb, err := db.Query("SELECT * FROM task WHERE id='" + strconv.Itoa(taskId) + "'")
	if err != nil {
		panic(err.Error())
	}
	task := Task{}

	for selDb.Next() {
		var id, done, userId int
		var title string

		err = selDb.Scan(&id, &userId, &title, &done)

		if err != nil {
			panic(err.Error())
		}
		task.Id = id
		task.UserId = userId
		task.Title = title
		task.Done = done
	}

	if task.Id == 0 {
		return Task{}
	}

	return task
}
func TaskCreate(c echo.Context) bool {
	db := db.GetInstance()

	task := Task{}
	err := c.Bind(&task)
	if err != nil {
		return false
	}

	insForm, err := db.Prepare("INSERT INTO task(user_id, title, done) VALUES(?,?,?)")
	if err != nil {
		return false
	}
	insForm.Exec(CurrentUser.Id, task.Title, 0)
	return true
}

func TaskDelete(taskId int) bool {
	db := db.GetInstance()

	delDb, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		return false
	}
	delDb.Exec(taskId)
	return true
}
func TaskUpdate(taskId int, c echo.Context) bool {
	db := db.GetInstance()
	task := Task{}
	err := c.Bind(&task)
	if err != nil {
		log.Println("error: ", err)
		return false
	}
	log.Println("task: ", task.Title)
	updDb, err := db.Prepare("UPDATE task SET title=?, done=? WHERE id=?")
	if err != nil {
		return false
	}

	updDb.Exec(task.Title, task.Done, taskId)
	return true
}
