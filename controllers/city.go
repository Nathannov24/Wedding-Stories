package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateCityController(c echo.Context) error {
	id := middlewares.ExtractTokenUserId(c)
	dataLogin, _ := database.GetUser(id)
	if dataLogin.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	err := database.InsertCity()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create city"))
}

func CreateNewCityController(c echo.Context) error {
	city := models.City{}
	c.Bind(&city)
	id := middlewares.ExtractTokenUserId(c)
	dataLogin, _ := database.GetUser(id)
	if dataLogin.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	respon, err := database.InsertNewCity(city)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create city "+respon.County))
}

func GetCityController(c echo.Context) error {
	cities, err := database.GetAllCity()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all cities", cities))
}
