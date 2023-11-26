package services

import (
	"ci_cd/features/coupons"
	"ci_cd/helper/jwt"
	"errors"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type CouponService struct {
	m coupons.Repository
}

func New(model coupons.Repository) coupons.Service {
	return &CouponService{
		m: model,
	}
}

func (cs *CouponService) TambahKupon(token *golangjwt.Token, newCoupon coupons.Coupon) (coupons.Coupon, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return coupons.Coupon{}, err
	}

	result, err := cs.m.InsertKupon(userID, newCoupon)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return coupons.Coupon{}, errors.New("kupon sudah ada pada sistem")
		}
		return coupons.Coupon{}, errors.New("terjadi kesalahan server")
	}

	return result, nil
}

func (cs *CouponService) GetKuponByUser(token *golangjwt.Token) ([]coupons.Coupon, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return nil, err
	}

	result, err := cs.m.ReadKuponByUser(userID)

	if err != nil {
		return nil, errors.New("terjadi kesalahan server")
	}

	return result, nil
}

func (cs *CouponService) GetKupon(page int64, pageSize int64) ([]coupons.Coupon, coupons.Pagination, error) {
	result, pagination, err := cs.m.ReadKupon(page, pageSize)

	if err != nil {
		return nil, coupons.Pagination{}, errors.New("terjadi kesalahan server")
	}

	return result, pagination, nil
}
