package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

func GetInvoiceAdmin() ([]models.PaymentInvoice, error) {
	var paymentinvoice []models.PaymentInvoice
	query := config.DB.Table("payments").Select("payments.id, payments.reservation_id, payments.url_photo, users.name, users.email, organizers.wo_name, packages.package_name, reservations.date, reservations.total_pax, reservations.status_payment, packages.price, reservations.total_price").
		Joins("join reservations on reservations.id = payments.reservation_id").Joins("join packages on packages.id = reservations.package_id").Joins("join users on users.id = reservations.user_id").Joins("join organizers on organizers.id = packages.organizer_id").
		Where("reservations.status_payment = 'unpaid' AND reservations.status_order = 'accepted' AND reservations.deleted_at is NULL").Find(&paymentinvoice)
	if query.Error != nil {
		return nil, query.Error
	}
	return paymentinvoice, nil
}

func ChangePaymentStatus(idReserve int) (interface{}, error) {
	var reserve models.Reservation
	query := config.DB.Find(&reserve, idReserve)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	updateQuery := config.DB.Model(&models.Reservation{}).Where("id = ?", idReserve).Update("status_payment", "paid")
	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return reserve, nil
}
