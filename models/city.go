package models

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT"`
	County    string `gorm:"type:varchar(255);not null" json:"county" form:"county"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type GetCityRespon struct {
	ID     int
	County string
}
