package controllers

import (
	"alta-wedding/lib/database"
	responses "alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

//register user
func RegisterUsersController(c echo.Context) error {
	var user models.User
	// REGEX
	var pattern string
	var matched bool
	// Bind all data from JSON
	c.Bind(&user)
	// Check data cannot be empty
	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("input data cannot be empty"))
	}
	// Check Name cannot less than 5 characters
	if len(user.Name) < 5 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("name cannot less than 5 characters"))
	}
	// Check Format Email
	EmailToLower := strings.ToLower(user.Email)
	pattern = `^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$`
	matched, _ = regexp.Match(pattern, []byte(EmailToLower))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email must contain email format"))
	}
	// Check Format Name
	pattern = `^\w(\w+ ?)*$`
	regex, _ := regexp.Compile(pattern)
	matched = regex.Match([]byte(user.Name))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid format name"))
	}
	// Check Length of Character of Password
	if len(user.Password) < 8 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("password cannot less than 8 characters"))
	}
	// Check Email is Exist
	duplicate, _ := database.GetUserByEmail(user.Email)
	if duplicate > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Email was used, try another email"))
	}
	// hash password bcrypt
	Password, _ := database.GeneratehashPassword(user.Password)
	user.Password = Password //replace old password
	user.Role = "User"
	_, err := database.RegisterUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create new user"))
}

//login users
func LoginUsersController(c echo.Context) error {
	login := models.UserLogin{}
	c.Bind(&login)
	users, err := database.LoginUsers(&login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if users == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid email or password"))
	}
	token, _ := middlewares.CreateToken(int(users.ID))
	return c.JSON(http.StatusCreated, responses.StatusSuccessLogin("login success", users.ID, token, users.Name, users.Role))
}

//get user
func GetUsersController(c echo.Context) error {
	loginuser := middlewares.ExtractTokenUserId(c)
	datauser, e := database.GetUser(loginuser)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	respon := models.Profile{
		ID:    datauser.ID,
		Name:  datauser.Name,
		Email: datauser.Email,
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get user", respon))
}

//update user
func UpdateUserController(c echo.Context) error {
	var user models.User
	// REGEX
	var pattern string
	var matched bool
	// Bind all data from JSON
	c.Bind(&user)
	updateuser := middlewares.ExtractTokenUserId(c)
	// Check data cannot be empty
	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("input data cannot be empty"))
	}
	// Check Name cannot less than 5 characters
	if len(user.Name) < 5 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("name cannot less than 5 characters"))
	}
	// Check Format Email
	pattern = `^\w+@\w+\.\w+$`
	matched, _ = regexp.Match(pattern, []byte(user.Email))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email must contain email format"))
	}
	// Check Format Name
	pattern = `^\w(\w+ ?)*$`
	regex, _ := regexp.Compile(pattern)
	matched = regex.Match([]byte(user.Name))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid format name"))
	}
	// Check Length of Character of Password
	if len(user.Password) < 8 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("password cannot less than 8 characters"))
	}
	// Check Email is Exist
	userdata, _ := database.GetUser(updateuser)
	if userdata.Email != user.Email {
		row, err := database.GetUserByEmail(user.Email)
		if row > 0 || err != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("Email was used, try another email"))
		}
	}
	// hash password bcrypt
	Password, _ := database.GeneratehashPassword(user.Password)
	user.Password = Password
	user.Role = "User"
	_, e := database.UpdateUser(updateuser, user)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success update user"))
}

//delete user by id
func DeleteUserController(c echo.Context) error {
	deleteuser := middlewares.ExtractTokenUserId(c)
	_, e := database.DeleteUser(deleteuser)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal service error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccess("success delete user"))
}

func GetUserControllersTest() echo.HandlerFunc {
	return GetUsersController
}

func UpdateUserControllersTest() echo.HandlerFunc {
	return UpdateUserController
}

func DeleteUserControllersTest() echo.HandlerFunc {
	return DeleteUserController
}
