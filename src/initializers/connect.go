package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/session-auth/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustConnect() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, name)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")

	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully migrated the database")

	return db
}
