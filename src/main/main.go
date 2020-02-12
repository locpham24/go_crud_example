package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"main/db"
)

type User struct {
	Id   int
	Name string `json:"name"`
}

type Task struct {
	Id     int
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
	Done   int    `json:"done"`
}

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
var jwtKey = []byte("ToDoApp")

func getAll(c echo.Context) error {
	user := c.Get("current_user").(User)

	task := Task{}
	res := []Task{}

	db := db.GetInstance()
	selDb, err := db.Query("SELECT * FROM task WHERE user_id='" + strconv.Itoa(user.Id) + "' ORDER BY id DESC")
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

	resJson, err := json.Marshal(res)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}
	return c.String(http.StatusOK, string(resJson))
}

func create(c echo.Context) error {
	user := c.Get("current_user").(User)
	log.Println("create route: User Name: ", user.Name, "User ID: ", user.Id)

	db := db.GetInstance()
	defer c.Request().Body.Close()
	defer db.Close()

	task := Task{}
	err := c.Bind(&task)
	if err != nil {
		panic(err.Error())
	}

	insForm, err := db.Prepare("INSERT INTO task(user_id, title, done) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(user.Id, task.Title, 0)

	return c.String(http.StatusOK, "ok")
}

func show(c echo.Context) error {
	//get current user
	user := c.Get("current_user").(User)

	db := db.GetInstance()
	defer c.Request().Body.Close()

	id := c.Param("id")

	selDb, err := db.Query("SELECT * FROM task WHERE id='" + id + "'")
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
		return c.String(http.StatusNotFound, "Not found")
	}

	if task.UserId != user.Id {
		return c.String(http.StatusUnauthorized, "You have no permission to view this task")
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can not parse to json")
	}
	return c.String(http.StatusOK, string(taskJson))
}

func update(c echo.Context) error {
	db := db.GetInstance()
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
	db := db.GetInstance()
	fmt.Println("address of db is", db)
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

//////////////////// middlewares ///////////////////////
func authWithJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		claims := token.Claims.(jwt.MapClaims)
		userId, err := strconv.Atoi(claims["jti"].(string))

		if err != nil {

		}

		curUser := User{
			Id:   userId,
			Name: claims["name"].(string),
		}
		c.Set("current_user", curUser)

		log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

		return next(c)
	}
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "admin" && password == "admin" {
		//Create token
		token, err := createJwtToken("admin", "1")
		if err != nil {
			log.Println("Error Creating Jwt Token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":   token,
			"message": "You were log in!",
		})
	}

	if username == "chris" && password == "123456" {
		//Create token
		token, err := createJwtToken("chris", "2")
		if err != nil {
			log.Println("Error Creating Jwt Token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":   token,
			"message": "You were log in!",
		})
	}

	return c.String(http.StatusUnauthorized, "Your username or password is not correct")
}

func createJwtToken(name string, id string) (string, error) {
	jwtClaims := JwtClaims{
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        id,
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims)
	token, err := rawToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func main() {
	log.Println("Server started on: localhost:8080")
	e := echo.New()

	taskGroup := e.Group("/task")

	taskGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    jwtKey,
	}))

	taskGroup.Use(authWithJwt)

	taskGroup.GET("", getAll)
	taskGroup.POST("", create)
	taskGroup.GET("/:id", show)
	taskGroup.PUT("/:id", update)
	taskGroup.DELETE("/:id", delete)

	e.GET("/login", login)
	e.Logger.Fatal(e.Start(":8080"))
}
