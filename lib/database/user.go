package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

func GetUser(userID int) (*models.User, error) {
	var userid models.User
	tx := config.DB.Where("id=?", userID).Find(&userid)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		return &userid, nil
	}
	return nil, nil
}

func GetUserByEmail(email string) (int64, error) {
	var usermail models.User
	tx := config.DB.Where("email = ?", email).Find(&usermail)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected > 0 {
		return tx.RowsAffected, nil
	}
	return 0, nil
}

func RegisterUser(user models.User) (interface{}, error) {
	var userregister models.User
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return userregister, nil
}

func LoginUsers(user *models.UserLogin) (*models.User, error) {
	var err error
	userpassword := models.User{}
	if err = config.DB.Where("email = ?", user.Email).Find(&userpassword).Error; err != nil {
		return nil, err
	}
	check := CheckPasswordHash(user.Password, userpassword.Password)
	if !check {
		return nil, nil
	}
	return &userpassword, nil
}

func UpdateUser(id int, User models.User) (models.User, error) {
	var user models.User
	if err := config.DB.Find(&user, id).Error; err != nil {
		return user, err
	}
	expass, _ := GeneratehashPassword(User.Password)
	user.Name = User.Name
	user.Email = User.Email
	user.Password = expass
	if err := config.DB.Model(&models.User{}).Where("id=?", id).Updates(User).Error; err != nil {
		return user, err
	}
	return user, nil
}

func DeleteUser(id int) (interface{}, error) {
	var userid models.User
	if err := config.DB.Where("id = ?", id).Delete(&userid).Error; err != nil {
		return nil, err
	}
	return userid, nil
}
