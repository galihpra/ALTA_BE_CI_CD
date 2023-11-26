package repository

import (
	"ci_cd/features/coupons/repository"
	"ci_cd/features/users"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Nama     string
	Username string
	Password string
	Coupons  []repository.CouponModel `gorm:"foreignKey:UserID"`
}

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) users.Repository {
	return &userQuery{
		db: db,
	}
}

func (uq *userQuery) InsertUser(newUser users.User) (users.User, error) {
	var inputDB = new(UserModel)
	inputDB.Username = newUser.Username
	inputDB.Nama = newUser.Nama
	inputDB.Password = newUser.Password

	if err := uq.db.Create(&inputDB).Error; err != nil {
		return users.User{}, err
	}

	newUser.ID = inputDB.ID

	return newUser, nil
}

func (uq *userQuery) Login(username string) (users.User, error) {
	var userData = new(UserModel)

	if err := uq.db.Where("username = ?", username).First(userData).Error; err != nil {
		return users.User{}, err
	}

	var result = new(users.User)
	result.ID = userData.ID
	result.Username = userData.Username
	result.Nama = userData.Nama
	result.Password = userData.Password

	return *result, nil
}
