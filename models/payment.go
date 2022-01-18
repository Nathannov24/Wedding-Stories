package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID             int `gorm:"primarykey"`
	Reservation_ID int `gorm:"type:int(100);NOT NULL; unique " json:"reservationid" form:"reservationid"`
	User_ID        int
	UrlPhoto       string `gorm:"type:varchar(255);NOT NULL" json:"invoice" form:"invoice"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type PaymentInvoice struct {
	ID             int
	Reservation_ID int
	Name           string
	Email          string
	WoName         string
	PackageName    string
	Date           string
	UrlPhoto       string
	Status_Payment string
	Total_Price    int
	Total_Pax      int
	Price          int
}
