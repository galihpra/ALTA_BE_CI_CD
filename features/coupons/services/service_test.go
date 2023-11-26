package services_test

import (
	"ci_cd/features/coupons"
	"ci_cd/features/coupons/mocks"
	"ci_cd/features/coupons/services"
	"ci_cd/helper/jwt"
	"errors"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

func TestInputCoupon(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	var repoData = coupons.Coupon{NamaProgram: "Undian Berhadiah", Link: "https://www.mozilla.org", Kode: "ub12345", Gambar: "gambar", UserID: userID}
	var falseData = coupons.Coupon{}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("InsertKupon", userID, repoData).Return(repoData, nil).Once()
		res, err := s.TambahKupon(token, repoData)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, repoData.NamaProgram, res.NamaProgram)
	})

	t.Run("Failed Case", func(t *testing.T) {
		repo.On("InsertKupon", userID, falseData).Return(falseData, errors.New("Database Error")).Once()
		res, err := s.TambahKupon(token, falseData)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, uint(0), res.UserID)
		assert.Equal(t, "", res.NamaProgram)
	})
}

func TestReadCoupon(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	dataResult := []coupons.Coupon{
		{ID: 1, NamaProgram: "Undian Berhadiah", Link: "https://www.mozilla.org", UserID: userID, Gambar: "image"},
		{ID: 2, NamaProgram: "Gebyar Berhadiah", Link: "https://www.mozilla.org", UserID: userID, Gambar: "image"},
	}

	paginations := coupons.Pagination{
		TotalRecords: 20,
		TotalPages:   2,
		NextPage:     2,
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("ReadKupon", int64(1), int64(10)).Return(dataResult, paginations, nil).Once()
		couponsResult, paginationResult, err := s.GetKupon(1, 10)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, userID, couponsResult[0].UserID)
		assert.Equal(t, paginations, paginationResult)
	})

	t.Run("Failed Case", func(t *testing.T) {
		paginations = coupons.Pagination{
			TotalRecords: 0,
			TotalPages:   0,
			NextPage:     0,
		}
		repo.On("ReadKupon", int64(1), int64(10)).Return(nil, coupons.Pagination{}, errors.New("repository error")).Once()
		couponsResult, _, err := s.GetKupon(1, 10)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, couponsResult)
		assert.Equal(t, paginations, coupons.Pagination{})
		assert.Equal(t, "terjadi kesalahan server", err.Error())
	})

}

func TestReadCouponByIdUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	dataResult := []coupons.Coupon{
		{ID: 1, NamaProgram: "Undian Berhadiah", Link: "https://www.mozilla.org", UserID: userID, Gambar: "image"},
		{ID: 2, NamaProgram: "Gebyar Berhadiah", Link: "https://www.mozilla.org", UserID: userID, Gambar: "image"},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("ReadKuponByUser", userID).Return(dataResult, nil).Once()
		couponsResult, err := s.GetKuponByUser(token)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, dataResult, couponsResult)
	})

	t.Run("Failed Case", func(t *testing.T) {
		repo.On("ReadKuponByUser", userID).Return(nil, errors.New("repository error")).Once()
		couponsResult, err := s.GetKuponByUser(token)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, couponsResult)
		assert.Equal(t, "terjadi kesalahan server", err.Error())
	})

}
