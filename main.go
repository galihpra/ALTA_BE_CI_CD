package main

import (
	"ci_cd/config"
	ch "ci_cd/features/coupons/handler"
	cr "ci_cd/features/coupons/repository"
	cs "ci_cd/features/coupons/services"
	uh "ci_cd/features/users/handler"
	ur "ci_cd/features/users/repository"
	us "ci_cd/features/users/services"
	"ci_cd/helper/enkrip"
	"ci_cd/routes"
	"ci_cd/utils/database"

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
