package routers

import (
	"alta-wedding/constants"
	"alta-wedding/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	//CORS
	e.Use(echoMid.CORSWithConfig(echoMid.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete},
	}))

	// Midleware Auth
	r := e.Group("")
	r.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	// ------------------------------------------------------------------
	// LOGIN & REGISTER USER
	// ------------------------------------------------------------------
	e.POST("/register/users", controllers.RegisterUsersController)
	e.POST("/login/users", controllers.LoginUsersController)
	// ------------------------------------------------------------------
	// USER ROUTER
	// ------------------------------------------------------------------
	r.GET("/users/profile", controllers.GetUsersController)
	r.PUT("/users/profile", controllers.UpdateUserController)
	r.DELETE("/users/profile", controllers.DeleteUserController)
	// ------------------------------------------------------------------
	// LOGIN & REGISTER ORGANIZER
	// ------------------------------------------------------------------
	e.POST("/register/organizer", controllers.CreateOrganizerController)
	e.POST("/login/organizer", controllers.LoginOrganizerController)
	e.GET("/organizer/profile/:id", controllers.GetProileOrganizerbyIDController)
	r.GET("/organizer/profile", controllers.GetProfileOrganizerController)
	r.PUT("/organizer/profile", controllers.UpdateOrganizerController)
	r.PUT("/organizer/profile/photo", controllers.UpdatePhotoOrganizerController)
	r.PUT("/organizer/profile/document", controllers.UpdateDocumentsOrganizerController)
	r.GET("/order/organizer/my", controllers.GetMyReservationListController)
	// ------------------------------------------------------------------
	// PACKAGE
	// ------------------------------------------------------------------
	r.POST("/package", controllers.InsertPackageController)
	e.GET("/package", controllers.GetAllPackageController)
	e.GET("/package/:id", controllers.GetPackageByIDController)
	r.GET("/package/my", controllers.GetMyPackageController)
	r.DELETE("/package/:id", controllers.DeletePackageController)
	r.PUT("/package/:id", controllers.UpdatePackageController)
	r.PUT("/package/photo/:id", controllers.UpdatePhotoPackageController)

	// ------------------------------------------------------------------
	// RESERVATION
	// ------------------------------------------------------------------
	r.POST("/reservation", controllers.CreateReservationController)
	r.PUT("/order/status/:id", controllers.AcceptDeclineController)
	r.GET("/order/users/my", controllers.GetReservationController)

	r.POST("/payment/invoice", controllers.PostPaymentController)
	r.GET("/payment/invoice", controllers.GetInvoiceController)

	// ------------------------------------------------------------------
	// ADMIN Authorize
	// ------------------------------------------------------------------
	r.POST("/cities", controllers.CreateCityController)
	r.POST("/cities/new", controllers.CreateNewCityController)
	e.GET("/cities", controllers.GetCityController)
	r.PUT("/payment/invoice", controllers.ChangePaymentStatusController)

	return e
}
