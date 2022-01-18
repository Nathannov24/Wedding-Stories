package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID         int    `gorm:"primarykey"`
	Package_ID int    `gorm:"primarykey" json:"package_id" form:"package_id"`
	Photo_Name string `gorm:"type:varchar(50);not null" json:"photo_name" form:"photo_name"`
	UrlPhoto   string `gorm:"type:longtext" json:"urlphoto" form:"urlphoto"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Get_Photo struct {
	Package_ID int
	Nama_Photo string
}

type EditPhoto struct {
	Photo_Name string
	UrlPhoto   string
}
