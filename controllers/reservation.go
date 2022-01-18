package controllers

import (
	"alta-wedding/lib/database"
	responses "alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller untuk memasukkan barang baru ke Reservation
func CreateReservationController(c echo.Context) error {
	Reservation := models.Reservation{}
	c.Bind(&Reservation)
	logged := middlewares.ExtractTokenUserId(c)                  // EXTRACT TOKEN LOGIN
	input, _ := database.GetPackagesByID(Reservation.Package_ID) // GET PACKAGE ID
	if input == nil {
		return c.JSON(http.StatusBadRequest, responses.ReservationFailed()) // WRONG INPUT
	}
	Reservation.User_ID = logged                            // TRANSFER DATA FROM TOKEN
	respon, err := database.CreateReservation(&Reservation) // CREATE RESERVATION
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error")) // DATABASE OR SERVER INTERNAL ERROR
	}
	if respon == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("already reserve")) // DATABASE ALREADY RESERVATION
	}
	database.AddTotalPrice(Reservation.Package_ID, respon.ID)         // PERKALIAN DI QUERY
	return c.JSON(http.StatusCreated, responses.ReservationSuccess()) // RESERVATION SUCCESS
}

func GetReservationController(c echo.Context) error {
	logged := middlewares.ExtractTokenUserId(c)
	input, e := database.GetReservation(logged)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success", input))
}
