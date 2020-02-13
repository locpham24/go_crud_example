package api

import (
	"encoding/json"
	"github.com/labstack/echo"
	"main/model"
	"net/http"
	"strconv"
)

func GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		tasks := model.Tasks{}

		tasks = model.TaskAll(model.CurrentUser.Id)

		resJson, err := json.Marshal(tasks)

		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}
		return c.String(http.StatusOK, string(resJson))
	}
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		success := model.TaskCreate(c)

		if success == false {
			return c.String(http.StatusInternalServerError, "Can not create task")
		}

		return c.String(http.StatusOK, "ok")
	}
}

func GetOne() echo.HandlerFunc {
	return func(c echo.Context) error {
		//get current user
		task := model.Task{}

		id, err := strconv.Atoi(c.Param("id"))

		task = model.TaskOne(id)

		if task.Id == 0 {
			return c.String(http.StatusNotFound, "Not found")
		}

		if task.UserId != model.CurrentUser.Id {
			return c.String(http.StatusUnauthorized, "You have no permission to view this task")
		}

		taskJson, err := json.Marshal(task)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Can not parse to json")
		}
		return c.String(http.StatusOK, string(taskJson))
	}
}

func Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {

		}

		success := model.TaskUpdate(id, c)

		if success == false {
			return c.String(http.StatusInternalServerError, "Can not create task")
		}

		return c.String(http.StatusOK, "ok")
	}
}

func Delete() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {

		}

		success := model.TaskDelete(id)

		if success == false {
			return c.String(http.StatusInternalServerError, "Can not delete task")
		}
		return c.String(http.StatusOK, "ok")
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		m := echo.Map{}
		if err := c.Bind(&m); err != nil {
			return err
		}
		username := m["username"]
		password := m["password"]

		token, err := model.ValidateUser(username, password)

		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token":   token,
			"message": "You were log in!",
		})
	}
}
