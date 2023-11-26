package handler

import (
	"ci_cd/features/users"
	"ci_cd/helper/jwt"
	"ci_cd/helper/responses"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	s users.Service
}

func New(s users.Service) users.Handler {
	return &userHandler{
		s: s,
	}
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		var inputProcess = new(users.User)
		inputProcess.Nama = input.Nama
		inputProcess.Password = input.Password
		inputProcess.Username = input.Username

		result, err := uh.s.Register(*inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar pada sistem"
			}

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		var response = new(RegisterResponse)
		response.Username = result.Username
		response.Nama = result.Nama
		response.ID = result.ID

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)
		// return c.JSON(http.StatusCreated, map[string]any{
		// 	"message": "success create data",
		// 	"data":    response,
		// })
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := uh.s.Login(input.Username, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data yang diinputkan tidak ditemukan",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var response = new(LoginResponse)
		response.Username = result.Username
		response.Nama = result.Nama
		response.ID = result.ID
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "login success",
			"data":    response,
		})
	}
}
