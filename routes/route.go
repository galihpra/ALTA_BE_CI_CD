package routes

import (
	"ci_cd/features/coupons"
	"ci_cd/features/users"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uh users.Handler, ch coupons.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	routeUser(e, uh)
	routeCoupon(e, ch)
}

func routeUser(e *echo.Echo, uh users.Handler) {
	e.POST("/users", uh.Register())
	e.POST("/login", uh.Login())
}

func routeCoupon(e *echo.Echo, ch coupons.Handler) {
	e.POST("/coupons", ch.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/coupons", ch.Read())
	e.GET("/couponsbyuser", ch.ReadByUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
