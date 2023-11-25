package main

import (
	"TugasDay23/config"
	ch "TugasDay23/features/coupons/handler"
	cr "TugasDay23/features/coupons/repository"
	cs "TugasDay23/features/coupons/services"
	uh "TugasDay23/features/users/handler"
	ur "TugasDay23/features/users/repository"
	us "TugasDay23/features/users/services"
	"TugasDay23/helper/enkrip"
	"TugasDay23/routes"
	"TugasDay23/utils/database"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	cfg := config.InitConfig()

	if cfg == nil {
		e.Logger.Fatal("tidak bisa start karena ENV error")
		return
	}

	db, err := database.InitMySQL(*cfg)

	if err != nil {
		e.Logger.Fatal("tidak bisa start karena DB error:", err.Error())
		return
	}

	db.AutoMigrate(&ur.UserModel{}, &cr.CouponModel{})

	userRepo := ur.New(db)
	hash := enkrip.New()
	userService := us.New(userRepo, hash)
	userHandler := uh.New(userService)

	couponRepo := cr.New(db)
	couponService := cs.New(couponRepo)
	couponHandler := ch.New(couponService)

	routes.InitRoute(e, userHandler, couponHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
