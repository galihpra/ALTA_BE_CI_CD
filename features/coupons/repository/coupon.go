package repository

import (
	"ci_cd/features/coupons"
	"log"

	"gorm.io/gorm"
)

type CouponModel struct {
	gorm.Model
	NamaProgram string
	Kode        string
	Link        string
	Gambar      string
	UserID      uint
}

type couponQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) coupons.Repository {
	return &couponQuery{
		db: db,
	}
}

func (cq *couponQuery) InsertKupon(userID uint, newCoupon coupons.Coupon) (coupons.Coupon, error) {
	var inputData = new(CouponModel)
	inputData.UserID = userID
	inputData.NamaProgram = newCoupon.NamaProgram
	inputData.Link = newCoupon.Link
	inputData.Kode = newCoupon.Kode
	inputData.Gambar = newCoupon.Gambar

	if err := cq.db.Create(&inputData).Error; err != nil {
		return coupons.Coupon{}, err
	}

	newCoupon.ID = inputData.ID

	return newCoupon, nil
}

func (cq *couponQuery) ReadKuponByUser(userID uint) ([]coupons.Coupon, error) {
	var couponModels []CouponModel

	if err := cq.db.Where("user_id = ?", userID).Find(&couponModels).Error; err != nil {
		return nil, err
	}

	var result []coupons.Coupon
	for _, model := range couponModels {
		result = append(result, coupons.Coupon{
			ID:          model.ID,
			NamaProgram: model.NamaProgram,
			Kode:        model.Kode,
			Link:        model.Link,
			Gambar:      model.Gambar,
			UserID:      model.UserID,
		})
	}

	return result, nil
}

func (cq *couponQuery) ReadKupon(page int64, pageSize int64) ([]coupons.Coupon, coupons.Pagination, error) {
	var couponModels []CouponModel

	var totalRecords int64
	cq.db.Model(&CouponModel{}).Count(&totalRecords)

	var offset = (page - 1) * pageSize
	if err := cq.db.Offset(int(offset)).Limit(int(pageSize)).Find(&couponModels).Error; err != nil {
		return nil, coupons.Pagination{}, err
	}

	var result []coupons.Coupon
	for _, model := range couponModels {
		result = append(result, coupons.Coupon{
			ID:          model.ID,
			NamaProgram: model.NamaProgram,
			Kode:        model.Kode,
			Link:        model.Link,
			Gambar:      model.Gambar,
			UserID:      model.UserID,
		})
	}

	var totalPages = (totalRecords + pageSize - 1) / pageSize

	var nextPage *int64
	if page < totalPages {
		nextPage = new(int64)
		*nextPage = page + 1
		log.Println("Next page:", *nextPage)
	}

	var prevPage *int64
	if page > 1 {
		prevPage = new(int64)
		*prevPage = page - 1
		log.Println("Previous page:", *prevPage)
	}

	var pagination = coupons.Pagination{
		TotalRecords: totalRecords,
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     *nextPage,
		PrevPage:     prevPage,
	}

	return result, pagination, nil
}
