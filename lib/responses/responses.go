package responses

import (
	"net/http"
)

// Fungsi untuk memberikan respon ketika controller gagal dijalankan
func StatusFailed(message string) map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "failed",
		"message": message,
	}
	return result
}

// Fungsi untuk memberikan respon ketika gagal upload photo
func StatusFailedDataPhoto(data interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"error":   true,
		"message": data,
	}
	return result
}

// Fungsi untuk memberikan respon ketika controller service error dijalankan
func StatusFailedInternal(message string, data interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "Unauthorized failed",
		"message": message,
		"data":    data,
	}
	return result
}

// Fungsi untuk memberikan respon ketika Authorisasi gagal
func StatusUnauthorized() map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "Unauthorized",
		"message": "Unauthorized Access",
	}
	return result
}

// Fungsi untuk memberikan respon ketika controller berhasil dijalankan
func StatusSuccess(message string) map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "success",
		"message": message,
	}
	return result
}

// Fungsi untuk memberikan respon ketika controller berhasil dijalankan dan menerima masukan data
func StatusSuccessData(message string, data interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
	return result
}

// function response false param
func FalseParamResponse() map[string]interface{} {
	result := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "False Param",
	}
	return result
}

// function response failed to reserve
func ReservationFailed() map[string]interface{} {
	result := map[string]interface{}{
		"status":  "failed",
		"message": "Failed to Reserve",
	}
	return result
}

// function response Success to reserve
func ReservationSuccess() map[string]interface{} {
	result := map[string]interface{}{
		"message": "Success to Reserve",
		"status":  "success",
	}
	return result
}

// function response wrong id
func WrongId() map[string]interface{} {
	result := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "Wrong Account",
	}
	return result
}

// function response Success to reserve
func SuccessCancelBook() map[string]interface{} {
	result := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success Cancel Reserve",
	}
	return result
}

// function response success to login with id display
func StatusSuccessLogin(message string, id, token, name, role interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"status":  "success",
		"message": message,
		"role":    role,
		"id":      id,
		"token":   token,
		"name":    name,
	}
	return result
}
