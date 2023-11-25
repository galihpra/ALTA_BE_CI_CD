package handler

import (
	"TugasDay23/features/coupons"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CouponHandler struct {
	s coupons.Service
}

func New(s coupons.Service) coupons.Handler {
	return &CouponHandler{
		s: s,
	}
}

func (cc *CouponHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(AddRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		var urlCloudinary = "cloudinary://254391669815624:nwV4h3433wh79UnqIog9ddzwvcY@dhxzinjxp"

		fileHeader, err := c.FormFile("gambar")

		validate := validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(input); err != nil {
			c.Echo().Logger.Error("Input error :", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
				"data":    nil,
			})
		}

		var inputProcess = new(coupons.Coupon)
		if err != nil {
			inputProcess.NamaProgram = input.NamaProgram
			inputProcess.Link = input.Link
			inputProcess.Kode = input.Kode
		} else {
			log.Println(fileHeader.Filename)

			file, _ := fileHeader.Open()

			var ctx = context.Background()

			cldService, _ := cloudinary.NewFromURL(urlCloudinary)
			resp, _ := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})
			log.Println(resp.SecureURL)

			inputProcess.NamaProgram = input.NamaProgram
			inputProcess.Link = input.Link
			inputProcess.Kode = input.Kode
			inputProcess.Gambar = resp.SecureURL
		}

		result, err := cc.s.TambahKupon(c.Get("user").(*gojwt.Token), *inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data yang diinputkan sudah terdaftar ada sistem",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(AddResponse)
		response.ID = result.ID
		response.NamaProgram = result.NamaProgram
		response.Kode = result.Kode
		response.Link = result.Link
		response.Gambar = result.Gambar

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success create data",
			"data":    response,
		})
	}
}

func (cc *CouponHandler) ReadByUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := cc.s.GetKuponByUser(c.Get("user").(*gojwt.Token))

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response []GetResponse
		for _, v := range result {
			var item = GetResponse{
				ID:          v.ID,
				NamaProgram: v.NamaProgram,
				Kode:        v.Kode,
				Link:        v.Link,
				Gambar:      v.Gambar,
				UserID:      v.UserID,
			}
			response = append(response, item)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success read data",
			"data":    response,
		})
	}
}

func (cc *CouponHandler) Read() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		result, pagination, err := cc.s.GetKupon(int64(page), int64(pageSize))

		if err != nil {
			c.Logger().Error("ERROR Get Data, explain:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response []GetResponse
		for _, v := range result {
			var item = GetResponse{
				ID:          v.ID,
				NamaProgram: v.NamaProgram,
				Kode:        v.Kode,
				Link:        v.Link,
				Gambar:      v.Gambar,
				UserID:      v.UserID,
			}
			response = append(response, item)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message":    "success read data",
			"data":       response,
			"pagination": pagination,
		})
	}
}
