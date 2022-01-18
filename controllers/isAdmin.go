package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReserveStatusBody struct {
	ID int `json:"reservationid" form:"reservationid"`
}

func GetInvoiceController(c echo.Context) error {
	invoiceadmin := middlewares.ExtractTokenUserId(c)
	datauser, e := database.GetUser(invoiceadmin)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if datauser.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	payment, err := database.GetInvoiceAdmin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}

	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get invoice", payment))
}

func ChangePaymentStatusController(c echo.Context) error {
	var input ReserveStatusBody
	c.Bind(&input)

	// Pengecekan user
	idToken := middlewares.ExtractTokenUserId(c)
	datauser, e := database.GetUser(idToken)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if datauser.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	// Mengubah status payment
	order, er := database.ChangePaymentStatus(input.ID)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to change payment status"))
	}
	if order == 0 {
		return c.JSON(http.StatusNotFound, responses.StatusFailed("reservation not found"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success to update payment status"))
}
