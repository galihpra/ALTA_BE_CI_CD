package users

import "github.com/labstack/echo/v4"

type User struct {
	ID       uint
	Nama     string
	Username string
	Password string
}

type Handler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
}

type Service interface {
	Register(newUser User) (User, error)
	Login(username string, password string) (User, error)
}

type Repository interface {
	InsertUser(newUser User) (User, error)
	Login(username string) (User, error)
}
