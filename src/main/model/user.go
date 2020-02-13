package model

import (
	"log"
)

type User struct {
	Id   int
	Name string `json:"name"`
}

var CurrentUser User

func ValidateUser(username interface{}, password interface{}) (string, error) {
	var err error
	if username == "admin" && password == "admin" {
		//Create token
		token, err := CreateJwtToken("admin", "1")
		if err != nil {
			log.Println("Error Creating Jwt Token:", err)
			return "", err
		}
		return token, nil
	}

	if username == "chris" && password == "123456" {
		//Create token
		token, err := CreateJwtToken("chris", "2")
		if err != nil {
			log.Println("Error Creating Jwt Token:", err)
			return "", err
		}
		return token, nil
	}

	return "", err
}
