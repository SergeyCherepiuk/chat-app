package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/chat-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresMustConnect() *gorm.DB {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DBNAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
