package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID             int     `gorm:"primarykey;AUTO_INCREMENT"`
	Package_ID     int     `json:"package_id" form:"package_id"`
	User_ID        int     `json:"user_id" form:"user_id"`
	Date           string  `gorm:"type:varchar(255)" json:"date" form:"date"`
	Additional     string  `gorm:"type:varchar(255)" json:"additional" form:"additional"`
	Total_Pax      int     `gorm:"type:int" json:"total_pax" form:"total_pax"`
	Total_Price    int     `gorm:"type:int;NOT NULL"`
	Status_Order   string  `gorm:"type:varchar(50); default:waiting" json:"status_order" form:"status_order"`
	Status_Payment string  `gorm:"type:varchar(50); default:unpaid" json:"status_payment" form:"status_payment"`
	Payment        Payment `gorm:"foreignkey:Reservation_ID;references:ID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type GetReservationRespon struct {
	ID             int
	Package_ID     int
	WoName         string
	PhoneNumber    string
	PackageName    string
	Address        string
	Email          string
	Date           string
	Additional     string
	Total_Pax      int
	Status_Order   string
	Status_Payment string
}
