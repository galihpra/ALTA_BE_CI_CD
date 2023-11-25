package coupons

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

type Coupon struct {
	ID          uint
	NamaProgram string
	Kode        string
	Link        string
	Gambar      string
	UserID      uint
}

type Pagination struct {
	TotalRecords int64  `json:"total_records"`
	CurrentPage  int64  `json:"current_page"`
	TotalPages   int64  `json:"total_pages"`
	NextPage     int64  `json:"next_page"`
	PrevPage     *int64 `json:"prev_page"`
}

type Handler interface {
	Add() echo.HandlerFunc
	ReadByUser() echo.HandlerFunc
	Read() echo.HandlerFunc
}

type Service interface {
	TambahKupon(token *jwt.Token, newCoupon Coupon) (Coupon, error)
	GetKuponByUser(token *jwt.Token) ([]Coupon, error)
	GetKupon(page int64, pageSize int64) ([]Coupon, Pagination, error)
}

type Repository interface {
	InsertKupon(userID uint, newCoupon Coupon) (Coupon, error)
	ReadKuponByUser(userID uint) ([]Coupon, error)
	ReadKupon(page int64, pageSize int64) ([]Coupon, Pagination, error)
}
