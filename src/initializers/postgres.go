package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/session-auth/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresMustConnect() *gorm.DB {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	name := os.Getenv("POSTGRES_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, name)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the Postgres database")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully migrated the Postgres database")

	return db
}
