package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Task struct {
	Id    int
	Title string `json:"title"`
	Done  int    `json:"done"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "go_crud_example"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getAll(c echo.Context) error {
	task := Task{}
	res := []Task{}

	db := dbConn()
	selDb, err := db.Query("SELECT * FROM task ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	for selDb.Next() {
		var id, done int
		var title string
		err = selDb.Scan(&id, &title, &done)
		if err != nil {
			panic(err.Error())
		}
		task.Id = id
		task.Title = title
		task.Done = done

		res = append(res, task)
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}
	return c.String(http.StatusOK, string(resJson))
}

func create(c echo.Context) error {
	db := dbConn()
	defer c.Request().Body.Close()
	defer db.Close()

	task := Task{}
	err := c.Bind(&task)
	if err != nil {
		panic(err.Error())
	}

	insForm, err := db.Prepare("INSERT INTO task(title,done) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(task.Title, 0)

	return c.String(http.StatusOK, "ok")
}

func show(c echo.Context) error {
	db := dbConn()
	defer c.Request().Body.Close()
	defer db.Close()

	id := c.Param("id")

	selDb, err := db.Query("SELECT * FROM task WHERE id='" + id + "'")
	if err != nil {
		panic(err.Error())
	}
	task := Task{}

	for selDb.Next() {
		var id, done int
		var title string

		err = selDb.Scan(&id, &title, &done)

		if err != nil {
			panic(err.Error())
		}
		task.Id = id
		task.Title = title
		task.Done = done
	}

	if task.Id == 0 {
		return c.String(http.StatusNotFound, "Not found")
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse to json")
	}
	return c.String(http.StatusOK, string(taskJson))
}

func update(c echo.Context) error {
	db := dbConn()
	defer c.Request().Body.Close()
	defer db.Close()

	id := c.Param("id")
	task := Task{}
	err := c.Bind(&task)

	updDb, err := db.Prepare("UPDATE task SET title=?, done=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	updDb.Exec(task.Title, task.Done, id)
	return c.String(http.StatusOK, "ok")
}

func delete(c echo.Context) error {
	db := dbConn()
	defer c.Request().Body.Close()
	defer db.Close()

	id := c.Param("id")

	delDb, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delDb.Exec(id)
	return c.String(http.StatusOK, "ok")
}

func main() {
	log.Println("Server started on: localhost:8080")
	e := echo.New()

	e.GET("/", getAll)
	e.POST("/", create)
	e.GET("/:id", show)
	e.PUT("/:id", update)
	e.DELETE("/:id", delete)
	e.Logger.Fatal(e.Start(":8080"))
}
