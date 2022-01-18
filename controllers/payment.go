package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

// Controller untuk memasukkan package baru
func PostPaymentController(c echo.Context) error {
	// Mendapatkan data package baru dari client
	input := models.Payment{}
	c.Bind(&input)
	user_id := middlewares.ExtractTokenUserId(c)
	input.User_ID = user_id
	// Menyimpan data barang baru menggunakan fungsi InsertPackage
	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("invoice")
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
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/heic" && filetype != "application/pdf" {
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
	invoicename := strings.ReplaceAll(strconv.Itoa(input.Reservation_ID), " ", "+")
	uploaded_file.Filename = fmt.Sprintf("invoice-%v-%v.%v", invoicename, formatted, extension)
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
	invoice := fmt.Sprintf("%v", u)
	input.UrlPhoto = invoice
	waiting, tx := database.PostPayment(input)
	if tx != nil {
		return c.JSON(http.StatusTeapot, responses.StatusFailed("payment was paid"))
	}
	if waiting == nil {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success sending invoice"))
}
