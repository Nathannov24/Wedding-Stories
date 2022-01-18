package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"encoding/json"
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
func InsertPackageController(c echo.Context) error {
	// Mendapatkan data package baru dari client
	input := models.Package{}
	c.Bind(&input)
	duplicate, _ := database.GetPackageByName(input.PackageName)
	if duplicate > 0 {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("package name was use, try input another package name"))
	}
	organizer_id := middlewares.ExtractTokenUserId(c)
	input.Organizer_ID = organizer_id
	// Menyimpan data barang baru menggunakan fungsi InsertPackage
	data, e := database.InsertPackage(input)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("failed to input package"))
	}

	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("urlphoto")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	lower := strings.ToLower(uploaded_file.Filename)
	fileExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "bmp": true}
	ext := strings.Split(lower, ".")
	extension := ext[len(ext)-1]
	if !fileExtensions[extension] {
		return c.JSON(http.StatusBadRequest, responses.StatusFailedDataPhoto("invalid type"))
	}

	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	packageName := strings.ReplaceAll(input.PackageName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("%v-%v.%v", packageName, formatted, extension)
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
	urlPhoto := fmt.Sprintf("%v", u)
	foto := models.Photo{
		Package_ID: data.ID,
		Photo_Name: packageName,
		UrlPhoto:   urlPhoto,
	}
	_, tx := database.InsertPhoto(foto)
	if tx != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}

	return c.JSON(http.StatusCreated, responses.StatusSuccess("success to input package"))
}

// Controller untuk mendapatkan seluruh data Packages
func GetAllPackageController(c echo.Context) error {
	// Mendapatkan data satu buku menggunakan fungsi GetPackages
	paket, e := database.GetPackages()
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch packages"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all packages", paket))
}

// Controller untuk mendapatkan seluruh data Packages by ID
func GetPackageByIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	// Mendapatkan data satu buku menggunakan fungsi GetPackages
	paket, e := database.GetPackagesByID(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch packages"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all packages by ID", paket))
}

func DeletePackageController(c echo.Context) error {
	// Mendapatkan id cart yang diingikan client
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}

	// Pengecekan apakah id package memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	getPackage, err := database.GetPackagesByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch package"))
	}
	getPackageJSON, err := json.Marshal(getPackage)
	if err != nil {
		panic(err)
	}

	var responsePackage models.Package
	json.Unmarshal([]byte(getPackageJSON), &responsePackage)

	if getPackage.Organizer_ID != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Unauthorized Access"))
	}

	// Mengapus data satu product menggunakan fungsi DeleteShoppingCart
	paket, e := database.DeletePackage(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to delete package"))
	}
	if paket == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("package id not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccess("success deleted package"))
}

// Update/Edit Package Function
func UpdatePackageController(c echo.Context) error {
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	paket := models.Package{}
	c.Bind(&paket)
	organizer_id := middlewares.ExtractTokenUserId(c)
	getPackage, _ := database.GetPackagesByID(id)
	if organizer_id != getPackage.Organizer_ID {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}

	// Edit into database
	database.UpdatePackage(id, paket)
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success edit data"))
}

// Update/Edit Profil Photo Package Function
func UpdatePhotoPackageController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	getPackageData, _ := database.GetPackagesByID(id)
	if organizer_id != getPackageData.Organizer_ID {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("urlphoto")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	lower := strings.ToLower(uploaded_file.Filename)
	fileExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "bmp": true}
	ext := strings.Split(lower, ".")
	extension := ext[len(ext)-1]
	if !fileExtensions[extension] {
		return c.JSON(http.StatusBadRequest, responses.StatusFailedDataPhoto("invalid type"))
	}

	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	organizerName := strings.ReplaceAll(getPackageData.PackageName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("%v-%v.%v", organizerName, formatted, extension)
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
	_, tx := database.UpdatePhotoPackage(urlLogo, id)
	if tx != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success upload photo"))
}
