package services_test

import (
	"ci_cd/features/users"
	"ci_cd/features/users/mocks"
	"ci_cd/features/users/services"
	enkMock "ci_cd/helper/enkrip/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := services.New(repo, enkrip)
	var inputData = users.User{Nama: "Galih", Username: "galihp30", Password: "gal123"}
	var repoData = users.User{Nama: "Galih", Username: "galihp30", Password: "some string"}
	var successReturnData = users.User{ID: uint(1), Nama: "Galih", Username: "galihp30"}
	var errReturnData = users.User{}
	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("some string", nil).Once()
		repo.On("InsertUser", repoData).Return(successReturnData, nil).Once()
		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "", res.Password)

	})

	t.Run("Hashing Error Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("", errors.New("hashing error")).Once()

		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Nama)
	})

	t.Run("Duplicate Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("some string", nil).Once()
		repo.On("InsertUser", repoData).Return(errReturnData, errors.New("duplicate entry")).Once()

		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Nama)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := services.New(repo, enkrip)

	var inputData = users.User{Username: "galihp30", Password: "gal123"}
	var successReturnData = users.User{ID: uint(1), Nama: "Galih", Username: "galihp30", Password: "hashed"}
	var errReturnData = users.User{}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("Login", inputData.Username).Return(successReturnData, nil).Once()
		enkrip.On("Compare", successReturnData.Password, inputData.Password).Return(nil).Once()
		res, err := s.Login(inputData.Username, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Galih", res.Nama)
		assert.Equal(t, inputData.Username, res.Username)
	})

	t.Run("Failed Not Found Case", func(t *testing.T) {
		repo.On("Login", inputData.Username).Return(errReturnData, errors.New("data tidak ditemukan")).Once()
		res, err := s.Login(inputData.Username, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, errReturnData, res)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
	})

	t.Run("Failed Incorrect Password Case", func(t *testing.T) {
		repo.On("Login", inputData.Username).Return(successReturnData, nil).Once()
		enkrip.On("Compare", successReturnData.Password, inputData.Password).Return(errors.New("terjadi kesalahan pada sistem")).Once()
		res, err := s.Login(inputData.Username, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, errReturnData, res)
		assert.Equal(t, "password yang diinputkan salah", err.Error())
	})
}
