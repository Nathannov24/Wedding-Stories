package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"alta-wedding/util"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

// Register Organizer Function
func CreateOrganizerController(c echo.Context) error {
	organizer := models.Organizer{}
	// Bind all data from JSON
	if err := c.Bind(&organizer); err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	// Check data cannot be empty
	if organizer.City == "" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("input data cannot be empty"))
	}
	// REGEX
	var pattern string
	var matched bool
	// Check Format Name
	pattern = `^(\w+ ?){4}$`
	regex, _ := regexp.Compile(pattern)
	matched = regex.Match([]byte(organizer.WoName))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("business name cannot less than 5 characters or invalid format"))
	}
	// Check Format Email
	emailLower := strings.ToLower(organizer.Email)
	pattern = `^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$`
	matched, _ = regexp.Match(pattern, []byte(emailLower))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email must contain email format"))
	}
	organizer.Email = emailLower
	// Check Format Password
	pattern = `^([a-zA-Z0-9()@:%_\+.~#?&//=\n"'\t\\;<>!$*-{}]+ ?){8}$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Password))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("password must contain password format and more than 8 characters"))
	}
	// Check Format Phone Number
	pattern = `^[0-9]{8,15}$`
	matched, _ = regexp.Match(pattern, []byte(organizer.PhoneNumber))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone number must contain phone number format"))
	}
	// Check Format Address
	pattern = `^[a-zA-Z]([a-zA-Z.0-9,]+ ?)*$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Address))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("address must be valid"))
	}
	// Check Organizer Email is Exist
	emailCheck, _ := database.CheckDatabase("email", organizer.Email)
	if emailCheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email was used, try another one"))
	}
	// Check Organizer Business name is Exist
	nameCheck, _ := database.CheckDatabase("wo_name", organizer.WoName)
	if nameCheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("business name was used, try another one"))
	}
	// Check Address
	_, _, Err := util.GetGeocodeLocations(organizer.Address)
	if Err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Address "+Err.Error()))
	}
	// Check Length of Character of PhoneNumber and Password
	if len(organizer.Password) < 8 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("password cannot less than 8 characters"))
	}
	// Check Phone number existing
	phonecheck, _ := database.CheckDatabase("phone_number", organizer.PhoneNumber)
	if phonecheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone number was used, try another one"))
	}
	// hash password bcrypt
	password, _ := database.GeneratehashPassword(organizer.Password)
	organizer.Password = password // replace old password to bcrypt password
	// Insert ALL data to Database
	_, e := database.InsertOrganizer(organizer)
	if e != nil {
		// Respon Failed
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	// Respon Success
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create new organizer"))
}

// Login Organizer Function
func LoginOrganizerController(c echo.Context) error {
	login := models.LoginRequestBody{}
	// Bind all data from JSON
	c.Bind(&login)
	emailLower := strings.ToLower(login.Email)
	login.Email = emailLower
	organizer, err := database.LoginOrganizer(login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if organizer == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid email or password"))
	}
	token, _ := middlewares.CreateToken(int(organizer.ID))
	return c.JSON(http.StatusCreated, responses.StatusSuccessLogin("login success", organizer.ID, token, organizer.WoName, "organizer"))
}

// Get Profile Organizer Function
func GetProfileOrganizerController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	respon, err := database.FindProfilOrganizer(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get organizer", respon))
}

// Get Profile Organizer by ID
func GetProileOrganizerbyIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	respon, err := database.FindProfilOrganizer(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get organizer", respon))
}

// Get my Package for Organizer
func GetMyPackageController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	mypackages, err := database.GetPackagesByToken(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get my packages", mypackages))
}

// Get My Reservation List From Users Order
func GetMyReservationListController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	mylistorder, err := database.GetListReservations(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get my list order", mylistorder))
}

// Fitur Accept/Decline Reservation
func AcceptDeclineController(c echo.Context) error {
	reservation_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	organizer_id := middlewares.ExtractTokenUserId(c)
	request := models.AcceptBody{}
	c.Bind(&request)
	orderstatus := strings.ToLower(request.Status_Order)
	// Check inputan harus accept atau decline
	pattern := `^\W*((?i)accept|decline(?-i))\W*$`
	matched, _ := regexp.Match(pattern, []byte(orderstatus))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("data must be accept/decline"))
	}
	request.Status_Order = orderstatus
	row, err := database.AcceptDecline(reservation_id, request.Status_Order, organizer_id)
	if err != nil || row < 1 {
		return c.JSON(http.StatusNotFound, responses.StatusFailed("Reservation Not Found"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success "+orderstatus+" reservation"))
}

// Update/Edit Profile Organizer Function
func UpdateOrganizerController(c echo.Context) error {
	organizer := models.Organizer{}
	// Bind all data from JSON
	if err := c.Bind(&organizer); err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	organizer_id := middlewares.ExtractTokenUserId(c)
	organizerData, _ := database.FindOrganizerById(organizer_id)
	// REGEX
	var pattern string
	var matched bool
	// Check Format Name
	pattern = `^(\w+ ?){4}$`
	regex, _ := regexp.Compile(pattern)
	matched = regex.Match([]byte(organizer.WoName))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid format name"))
	}
	// Check Format Email
	emailLower := strings.ToLower(organizer.Email)
	pattern = `^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Email))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email must contain email format"))
	}
	organizer.Email = emailLower
	// Check Password
	if organizer.Password == "" {
		organizer.Password = organizerData.Password
	} else {
		// Check Format Password
		pattern = `^([a-zA-Z0-9()@:%_\+.~#?&//=\n"'\t\\;<>!$*-{}]+ ?){8}$`
		matched, _ = regexp.Match(pattern, []byte(organizer.Password))
		if !matched {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("password must contain password format and more than 8 characters"))
		}
		// hash password bcrypt
		password, _ := database.GeneratehashPassword(organizer.Password)
		organizer.Password = password // replace old password to bcrypt password
	}
	// Check Format Phone Number
	pattern = `^[0-9]{8,15}$`
	matched, _ = regexp.Match(pattern, []byte(organizer.PhoneNumber))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone number must contain phone number format"))
	}
	// Check Format Web URL
	pattern = `^[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`
	matched, _ = regexp.Match(pattern, []byte(organizer.WebUrl))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("url web must contain url format"))
	}
	// Check Format About
	pattern = `^\w([a-zA-Z0-9()@:%_\+.~#?&//=\n"'\t\\;<>!*-{}]+ ?)*$`
	matched, _ = regexp.Match(pattern, []byte(organizer.About))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("description cannot be empty"))
	}
	// Check Format Address
	pattern = `^[a-zA-Z]([a-zA-Z.0-9,]+ ?)*$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Address))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Address must be valid"))
	}
	// Check Email Organizer is Exist
	if organizer.Email != organizerData.Email {
		row, err := database.CheckDatabase("email", organizer.Email)
		if row > 0 || err != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("email was used, try another one"))
		}
	}
	// Check Business name Organizer is Exist
	if organizer.WoName != organizerData.WoName {
		row, err := database.CheckDatabase("wo_name", organizer.WoName)
		if row > 0 || err != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("business name was used, try another one"))
		}
	}
	// Check Phone number existing
	if organizer.PhoneNumber != organizerData.PhoneNumber {
		phonecheck, er := database.CheckDatabase("phone_number", organizer.PhoneNumber)
		if phonecheck > 0 || er != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone number was used, try another one"))
		}
	}
	// Check Address Valid Apa Enggak
	_, _, Err := util.GetGeocodeLocations(organizer.Address)
	if Err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Address "+Err.Error()))
	}
	// Edit into database
	_, err := database.EditOrganizer(organizer, organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success edit data"))
}

// Update/Edit Profil Photo Organizer Function
func UpdatePhotoOrganizerController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	dataWo, _ := database.FindOrganizerById(organizer_id)
	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("logo")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	buff := make([]byte, 512)
	_, err = f.Read(buff)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	// Cek Ekstension Type must be JPEG or PNG
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("The provided file format is not allowed. Please upload a JPEG or PNG image"))
	}
	// Return the pointer back to the start of the file
	_, er := f.Seek(0, io.SeekStart)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(er.Error()))
	}
	if uploaded_file.Size > MAX_UPLOAD_SIZE*5 {
		return c.JSON(http.StatusBadGateway, responses.StatusFailed("The uploaded file is too big. Please choose an file that's less than 5MB in size"))
	}
	ext := strings.Split(uploaded_file.Filename, ".")
	extension := ext[len(ext)-1]
	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	organizerName := strings.ReplaceAll(dataWo.WoName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("Photo-%v-%v.%v", organizerName, formatted, extension)
	sw := storageClient.Bucket(bucket).Object(uploaded_file.Filename).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	if err := sw.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	u, err := url.Parse("https://storage.googleapis.com/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	// Insert URL
	urlLogo := fmt.Sprintf("%v", u)
	_, e := database.EditPhotoOrganizer(urlLogo, organizer_id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success upload photo"))
}

// Insert Document Organizer Function
func UpdateDocumentsOrganizerController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	dataWo, _ := database.FindOrganizerById(organizer_id)
	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	buff := make([]byte, 512)
	_, err = f.Read(buff)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	// Cek Ekstension Type must be JPEG or PNG
	filetype := http.DetectContentType(buff)
	if filetype != "application/pdf" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("The provided file format is not allowed. Please upload a PDF document"))
	}
	// Return the pointer back to the start of the file
	_, er := f.Seek(0, io.SeekStart)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(er.Error()))
	}
	if uploaded_file.Size > MAX_UPLOAD_SIZE*5 {
		return c.JSON(http.StatusBadGateway, responses.StatusFailed("The uploaded file is too big. Please choose an file that's less than 5MB in size"))
	}
	ext := strings.Split(uploaded_file.Filename, ".")
	extension := ext[len(ext)-1]
	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	organizerName := strings.ReplaceAll(dataWo.WoName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("Document-%v-%v.%v", organizerName, formatted, extension)
	sw := storageClient.Bucket(bucket).Object(uploaded_file.Filename).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	if err := sw.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	u, err := url.Parse("https://storage.googleapis.com/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	// Insert URL
	urlDocument := fmt.Sprintf("%v", u)
	_, e := database.EditDocumentOrganizer(urlDocument, organizer_id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	database.UpdateStatusWO(organizer_id)
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success upload document"))
}

//------------------------------------------------------
//>>>>>>>>>>>>>>>>> FOR UNIT TESTING <<<<<<<<<<<<<<<<<<<
//------------------------------------------------------
// Testing Get Profile Organizer
func GetProfileOrganizerControllerTest() echo.HandlerFunc {
	return GetProfileOrganizerController
}

// Testing Get My Reservation
func GetMyReservationListControllerTest() echo.HandlerFunc {
	return GetMyReservationListController
}

// Testing Get My Packages
func GetMyPackageControllerTest() echo.HandlerFunc {
	return GetMyPackageController
}

// Testing Accept/Decline Feature
func AcceptDeclineControllerTest() echo.HandlerFunc {
	return AcceptDeclineController
}

// Testing Accept/Decline Feature
func UpdateOrganizerControllerTest() echo.HandlerFunc {
	return UpdateOrganizerController
}
