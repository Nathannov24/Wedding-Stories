package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int       `gorm:"primarykey"`
	Name      string    `gorm:"type:varchar(255)" json:"name" form:"name"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email" form:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	Role      string    `gorm:"type:varchar(20)" json:"role" form:"role"`
	Payments  []Payment `gorm:"foreignKey:User_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type Profile struct {
	ID    int
	Name  string
	Email string
}
