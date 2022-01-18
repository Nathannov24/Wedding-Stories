package controllers

import (
	"alta-wedding/config"
	"alta-wedding/constants"
	"alta-wedding/lib/database"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

// Fungsi untuk menginisialisasi koneksi ke database test
func InitEchoTestUser() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

// Struct yang digunakan ketika test request success, dapat menampung banyak data
type UsersResponseSuccess struct {
	Status  string
	Message string
	Data    models.User
}

type GetRespon struct {
	Status  string
	Message string
	ID      int
	Name    string
	Email   string
	Data    models.User
}

// Struct yang digunakan ketika test request failed
type ResponFailed struct {
	Status  string
	Message string
}

var (
	mock_user_data = models.User{
		Name:     "armuh",
		Email:    "armuh@gmail.com",
		Password: "yourpasswd",
	}
)

type EditUserRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

var logindata = models.UserLogin{
	Email:    "armuh@gmail.com",
	Password: "yourpasswd",
}

type ResponSuccessLogin struct {
	Status  string `json:"status" form:"status"`
	Message string `json:"message" form:"message"`
	ID      int    `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Role    string `json:"role" form:"role"`
	Token   string `json:"token" form:"token"`
}

var expass string

// Fungsi untuk memasukkan data user test ke dalam database
func InsertMockUserDataToDB() error {
	expass, _ = database.GeneratehashPassword(mock_user_data.Password)
	mock_user_data.Password = expass
	var err error
	if err = config.DB.Save(&mock_user_data).Error; err != nil {
		return err
	}
	return nil
}

//test login success done
func TestLoginGetUserControllers(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockUserDataToDB()
	body, error := json.Marshal(logindata)
	if error != nil {
		t.Error(t, error, "error marshal")
	}
	// send data using request body with HTTP method POST
	req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/login/users")

	if assert.NoError(t, LoginUsersController(context)) {
		res_body := res.Body.String()
		var User ResponSuccessLogin
		err := json.Unmarshal([]byte(res_body), &User)
		if err != nil {
			assert.Error(t, err, "error marshal")
		}

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, "success", User.Status)
		assert.Equal(t, "login success", User.Message)

	}
}

//test login error done
func TestLoginUserFailed(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockUserDataToDB()

	t.Run("TestLoginUser_InvalidInput", func(t *testing.T) {
		logininfo, err := json.Marshal(models.UserLogin{Email: "fian@gmail.com", Password: "admins"})
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(logininfo))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/users")
		if assert.NoError(t, LoginUsersController(contex)) {
			bodyResponses := rec.Body.String()
			var User ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &User)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", User.Status)
			assert.Equal(t, "invalid email or password", User.Message)
		}
	})
	t.Run("TestLoginUser_ErrorDB", func(t *testing.T) {
		datalogin, err := json.Marshal(logindata)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(datalogin))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/users")
		config.DB.Migrator().DropTable(&models.User{})
		if assert.NoError(t, LoginUsersController(contex)) {
			bodyResponses := rec.Body.String()
			var User ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &User)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "failed", User.Status)
			assert.Equal(t, "internal server error", User.Message)
		}
	})
}

//test register user success done
func TestRegisterUserController(t *testing.T) {
	e := InitEchoTestUser()
	body, err := json.Marshal(mock_user_data)
	if err != nil {
		t.Error(t, err, "error marshal")
	}
	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/register/users")

	if assert.NoError(t, RegisterUsersController(c)) {
		bodyuser := rec.Body.String()
		var user UsersResponseSuccess
		err := json.Unmarshal([]byte(bodyuser), &user)
		if err != nil {
			assert.Error(t, err, "error marshal")
		}
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "success", user.Status)
		assert.Equal(t, "success create new user", user.Message)
	}
}

//test register user error done
func TestRegisterUserControllerError(t *testing.T) {
	e := InitEchoTestUser()
	t.Run("TestRegisterEmpty", func(t *testing.T) {
		type Login struct {
			Name     string
			Password string
		}
		var empty Login
		body, err := json.Marshal(empty)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "input data cannot be empty", user.Message)
		}
	})
	t.Run("TestRegisterNameLess", func(t *testing.T) {
		mock_user_data.Name = "Arif"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "name cannot less than 5 characters", user.Message)
		}
	})
	t.Run("TestRegisterEmailWasUsed", func(t *testing.T) {
		InsertMockUserDataToDB()
		mock_user_data.Email = "armuh@gmail.com"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "Email was used, try another email", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatName", func(t *testing.T) {
		mock_user_data.Name = "      armuh"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "invalid format name", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatEmail", func(t *testing.T) {
		mock_user_data.Email = "#armuh@gmail.com"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "email must contain email format", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatPassword", func(t *testing.T) {
		mock_user_data.Password = "12345"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "password cannot less than 8 characters", user.Message)
		}
	})
	t.Run("TestRegisterBadRequest", func(t *testing.T) {
		config.DB.Migrator().DropTable(models.User{})
		mock_user_data.Email = "armuh@gmail.com"
		mock_user_data.Password = "yourpasswd"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "bad request", user.Message)
		}
	})
}

//test get user success done
func TestGetUserControllerSuccess(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockUserDataToDB()
	var UserDetail models.User
	tx := config.DB.Where("email=? AND password=?", logindata.Email, expass).First(&UserDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(UserDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/profile")
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserControllersTest())(c)

	var user UsersResponseSuccess
	bodyuser := rec.Body.String()
	json.Unmarshal([]byte(bodyuser), &user)
	t.Run("GET/users/profile", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", user.Status)
		assert.Equal(t, "success get user", user.Message)
		assert.Equal(t, "armuh", user.Data.Name)
		assert.Equal(t, "armuh@gmail.com", user.Data.Email)
	})
}

//test get user error done
func TestGetUserControllerError(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockUserDataToDB()
	var UserDetail models.User
	tx := config.DB.Where("email=? AND password=?", logindata.Email, expass).First(&UserDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(UserDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/profile")
	config.DB.Migrator().DropTable(&models.User{})
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserControllersTest())(c)

	var user UsersResponseSuccess
	bodyuser := rec.Body.String()
	json.Unmarshal([]byte(bodyuser), &user)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "failed", user.Status)
	assert.Equal(t, "internal server error", user.Message)
}

//test update user success done
func TestUpdateUserControllerSucces(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockUserDataToDB()
	var UserDetail models.User
	tx := config.DB.Where("email=? AND password=?", logindata.Email, expass).First(&UserDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(UserDetail.ID))
	if err != nil {
		panic(err)
	}
	var update = models.User{
		Name:     "iniupdate",
		Email:    "iniupdate@gmail.com",
		Password: "qwertyuiop",
		Role:     "user",
	}
	body, err := json.Marshal(update)
	if err != nil {
		t.Error(t, err, "error marshal")
	}
	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/profile")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateUserControllersTest())(c)
	bodyuser := rec.Body.String()
	var user UsersResponseSuccess
	json.Unmarshal([]byte(bodyuser), &user)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "success", user.Status)
	assert.Equal(t, "success update user", user.Message)
}

/*

//test update user error done
func TestUpdateUserControllerError(t *testing.T) {
	e := InitEchoTestUser()
	t.Run("TestRegisterEmpty", func(t *testing.T) {
		type Login struct {
			Name     string
			Password string
		}
		var empty Login
		body, err := json.Marshal(empty)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, UpdateUserControllersTest(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "input data cannot be empty", user.Message)
		}
	})
	t.Run("TestRegisterNameLess", func(t *testing.T) {
		mock_user_data.Name = "Arif"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "name cannot less than 5 characters", user.Message)
		}
	})
	t.Run("TestRegisterEmailWasUsed", func(t *testing.T) {
		InsertMockUserDataToDB()
		mock_user_data.Email = "armuh@gmail.com"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// config.DB.Migrator().DropTable(models.User{})

		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "Email was used, try another email", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatName", func(t *testing.T) {
		mock_user_data.Name = "      armuh"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// config.DB.Migrator().DropTable(models.User{})

		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "invalid format name", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatEmail", func(t *testing.T) {
		mock_user_data.Email = "#armuh@gmail.com"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// config.DB.Migrator().DropTable(models.User{})

		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "email must contain email format", user.Message)
		}
	})
	t.Run("TestRegisterInvalidFormatPassword", func(t *testing.T) {
		mock_user_data.Password = "12345"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// config.DB.Migrator().DropTable(models.User{})

		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "password cannot less than 8 characters", user.Message)
		}
	})
	t.Run("TestRegisterBadRequest", func(t *testing.T) {
		config.DB.Migrator().DropTable(models.User{})
		mock_user_data.Email = "armuh@gmail.com"
		mock_user_data.Password = "yourpasswd"
		body, err := json.Marshal(mock_user_data)
		if err != nil {
			t.Error(t, err, "error marshal")
		}

		//send data using request body with HTTP Method POST
		req := httptest.NewRequest(http.MethodPost, "/register/users", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, RegisterUsersController(c)) {
			bodyResponses := rec.Body.String()
			var user ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &user)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", user.Status)
			assert.Equal(t, "bad request", user.Message)
		}
	})
}


/*

//--------------------------------------------------------
//test update user


//test update user
func TestUpdateUserControllerParam(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed Update User",
		path:       "/user/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("#")

	if assert.NoError(t, UpdateUserController(context)) {

		var response ResponFailed
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "False Param", response.Message)
	}

}

//test update user error
func TestUpdateUserControllerError(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed Update User",
		path:       "/user/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(&models.User{})
	req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, UpdateUserController(context)) {

		var response ResponFailed
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "Bad Request", response.Message)
	}

}

//--------------------------------------------------------

//test delete user
func TestDeleteUserController(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Success Delete User",
		path:       "/users/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodDelete, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/users/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")
	if assert.NoError(t, DeleteUserController(contex)) {
		var user ResponFailed
		body := res.Body.String()
		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error unmarshal")
		}
		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "Success Delete User", user.Message)
	}
}

//test delete user error
func TestDeleteUserControllerError(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed to Create User",
		path:       "/users/:id",
		expectCode: http.StatusBadRequest,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(&models.User{})
	req := httptest.NewRequest(http.MethodDelete, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath(testCases.path)
	c.SetParamNames("id")
	c.SetParamValues("1")

	//send data using request body with HTTP Method DELETE
	if assert.NoError(t, DeleteUserController(c)) {
		bodyResponses := res.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "Failed to delete user", user.Message)
	}

}

//-----------------------------------------
//tes login error
func TestLoginUsersControllersError(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(models.User{})
	req := httptest.NewRequest(http.MethodPost, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)

	if assert.NoError(t, LoginUsersController(context)) {

		var user UserResponse
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		// assert.Equal(t, testCases.expectStatus, res.Code)
		assert.Equal(t, "Login failed", user.Message)

	}
}

*/
