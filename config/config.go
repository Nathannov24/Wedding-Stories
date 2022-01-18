package config

import (
	"alta-wedding/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConfig() map[string]string {
	conf, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conf
}

var DB *gorm.DB

// Initial Database
func InitDB() {
	dbconfig := GetConfig()
	// Sesuaikan dengan database kalian
	connect := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		dbconfig["DB_USERNAME"],
		dbconfig["DB_PASSWORD"],
		dbconfig["DB_HOST"],
		dbconfig["DB_PORT"],
		dbconfig["DB_NAME"])

	var err error
	DB, err = gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	InitalMigration()
}

// Function Initial Migration (Tabel)
func InitalMigration() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Organizer{})
	DB.AutoMigrate(&models.Package{})
	DB.AutoMigrate(&models.Photo{})
	DB.AutoMigrate(&models.Reservation{})
	DB.AutoMigrate(&models.Payment{})
	DB.AutoMigrate(&models.City{})
}

// Initia Database Unit Testing
func InitDBTest() {
	dbconfig := map[string]string{
		"DB_USERNAME":  "root",
		"DB_PASSWORD":  "yourpasswd",
		"DB_HOST":      "localhost",
		"DB_PORT":      "3306",
		"DB_NAME_TEST": "alta_wedding_test"}
	// Sesuaikan dengan database kalian
	connect := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		dbconfig["DB_USERNAME"],
		dbconfig["DB_PASSWORD"],
		dbconfig["DB_HOST"],
		dbconfig["DB_PORT"],
		dbconfig["DB_NAME_TEST"])

	// Sesuaikan dengan database kalian
	var err error
	DB, err = gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	InitalMigrationTest()
}

// Function Initial Migration (Tabel)
func InitalMigrationTest() {
	DB.Migrator().DropTable(&models.Organizer{})
	DB.AutoMigrate(&models.Organizer{})
	DB.Migrator().DropTable(&models.User{})
	DB.AutoMigrate(&models.User{})
	DB.Migrator().DropTable(&models.Package{})
	DB.AutoMigrate(&models.Package{})
	DB.Migrator().DropTable(&models.Photo{})
	DB.AutoMigrate(&models.Photo{})
	DB.Migrator().DropTable(&models.User{})
	DB.AutoMigrate(&models.User{})
	DB.Migrator().DropTable(&models.Reservation{})
	DB.AutoMigrate(&models.Reservation{})
}
