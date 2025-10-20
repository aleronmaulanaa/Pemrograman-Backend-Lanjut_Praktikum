package repository

import (
	"praktikum4-crud/app/model"
	"praktikum4-crud/database"
)

func GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	err := database.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username=$1", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
