package services

import (
	"ci_cd/features/users"
	"ci_cd/helper/enkrip"
	"errors"
	"strings"
)

type userService struct {
	repo users.Repository
	h    enkrip.HashInterface
}

func New(r users.Repository, h enkrip.HashInterface) users.Service {
	return &userService{
		repo: r,
		h:    h,
	}
}

func (us *userService) Register(newUser users.User) (users.User, error) {
	// validasi

	// enkripsi password
	enkPassword, err := us.h.HashPassword(newUser.Password)

	if err != nil {
		return users.User{}, errors.New("terjadi masalah saat memproses data")
	}

	newUser.Password = enkPassword

	result, err := us.repo.InsertUser(newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return users.User{}, errors.New("data yang dimasukkan sudah terdaftar")
		}
		return users.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}
func (us *userService) Login(hp string, password string) (users.User, error) {
	result, err := us.repo.Login(hp)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return users.User{}, errors.New("data tidak ditemukan")
		}
		return users.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	err = us.h.Compare(result.Password, password)
	if err != nil {
		return users.User{}, errors.New("password yang diinputkan salah")
	}

	return result, nil
}
