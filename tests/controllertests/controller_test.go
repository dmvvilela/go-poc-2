package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dmvvilela/go-poc-2/api/controllers"
	"github.com/dmvvilela/go-poc-2/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var contactInstance = models.Contact{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshContactTable() error {
	err := server.DB.DropTableIfExists(&models.Contact{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Contact{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneContact() (models.Contact, error) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contact := models.Contact{
		Name:  "Pet",
		Email: "pet@gmail.com",
	}

	err = server.DB.Model(&models.Contact{}).Create(&contact).Error
	if err != nil {
		return models.Contact{}, err
	}

	return contact, nil
}

func seedContacts() ([]models.Contact, error) {
	var err error
	if err != nil {
		return nil, err
	}

	contacts := []models.Contact{
		{
			Name:  "Steven victor",
			Email: "steven@gmail.com",
		},
		{
			Name:  "Kenny Morris",
			Email: "kenny@gmail.com",
		},
	}

	for i := range contacts {
		err := server.DB.Model(&models.Contact{}).Create(&contacts[i]).Error
		if err != nil {
			return []models.Contact{}, err
		}
	}

	return contacts, nil
}
