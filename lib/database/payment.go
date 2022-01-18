package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

//
func PostPayment(Payment models.Payment) (*models.Payment, error) {
	query := config.DB.Where("id = ? AND status_order = 'accepted' AND status_payment = 'unpaid' AND user_id = ? ", Payment.Reservation_ID, Payment.User_ID).Find(
		&models.Reservation{})
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected > 0 {
		if err := config.DB.Create(&Payment).Error; err != nil {
			return nil, err
		}
		return &Payment, nil
	}
	return nil, nil
}
