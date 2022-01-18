package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

func InsertCity() error {
	city := []models.City{
		{County: "Jakarta"},
		{County: "Surabaya"},
		{County: "Bandung"},
		{County: "Bekasi"},
		{County: "Yogyakarta"},
		{County: "Tangerang"},
		{County: "Banten"},
		{County: "Semarang"},
		{County: "Bogor"},
		{County: "Malang"},
		{County: "Makassar"},
		{County: "Palu"},
		{County: "Manado"},
		{County: "Balikpapan"},
		{County: "Samarinda"},
		{County: "Singkawang"},
		{County: "Batam"},
		{County: "Pekanbaru"},
		{County: "Padang"},
		{County: "Lhoksemauwe"},
		{County: "Batam"},
		{County: "Denpasar"},
	}
	if err := config.DB.Create(&city).Error; err != nil {
		return err
	}
	return nil
}

func InsertNewCity(city models.City) (*models.City, error) {
	if err := config.DB.Create(&city).Error; err != nil {
		return nil, err
	}
	return &city, nil
}

func GetAllCity() ([]models.GetCityRespon, error) {
	cities := []models.GetCityRespon{}
	if err := config.DB.Table("cities").Select("cities.id, cities.county").Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}
