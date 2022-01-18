package models

import (
	"time"

	"gorm.io/gorm"
)

type Organizer struct {
	ID          int       `gorm:"primarykey;AUTO_INCREMENT"`
	WoName      string    `gorm:"type:varchar(255);not null" json:"woname" form:"woname"`
	Email       string    `gorm:"type:varchar(100);unique;not null" json:"email" form:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phonenumber" form:"phonenumber"`
	About       string    `gorm:"type:longtext" json:"about" form:"about"`
	City        string    `gorm:"type:varchar(255)" json:"city" form:"city"`
	Address     string    `gorm:"type:varchar(255)" json:"address" form:"address"`
	WebUrl      string    `gorm:"type:varchar(255)" json:"weburl" form:"weburl"`
	Proof       string    `gorm:"type:varchar(255)" json:"proof" form:"proof"`
	Logo        string    `gorm:"type:varchar(255)" json:"logo" form:"logo"`
	Status      string    `gorm:"type:varchar(255); default:Not Activated" json:"status" form:"status"`
	Package     []Package `gorm:"foreignKey:Organizer_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type PostRequestBody struct {
	WoName   string `json:"woname" form:"woname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	City     string `json:"city" form:"city"`
	Address  string `json:"address" form:"address"`
}

type LoginRequestBody struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type PostPhoto struct {
	Url string `json:"logo" form:"logo"`
}

type ProfileRespon struct {
	WoName      string `json:"woname" form:"woname"`
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phonenumber" form:"phonenumber"`
	About       string `json:"about" form:"about"`
	WebUrl      string `json:"weburl" form:"weburl"`
	Status      string `json:"status" form:"status"`
	Logo        string `json:"logo" form:"logo"`
	City        string `json:"city" form:"city"`
	Address     string `json:"address" form:"address"`
	Proof       string `json:"proof" form:"proof"`
}

type ReservationListRespon struct {
	ID             int
	Package_ID     int
	PackageName    string
	Name           string
	Date           string
	Additional     string
	Total_Pax      int
	Status_Order   string
	Status_Payment string
}

type AcceptBody struct {
	Status_Order string `json:"status_order" form:"status_order"`
}
