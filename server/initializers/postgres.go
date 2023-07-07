package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/session-auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresMustConnect() *gorm.DB {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	name := os.Getenv("POSTGRES_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s", user, password, host, port)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the Postgres database")

	createDatabaseCommand := fmt.Sprintf("SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')", name, name)
	r := db.Exec(createDatabaseCommand)
	if r.Error != nil {
		log.Fatal(r.Error)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully migrated the Postgres database")

	return db
}
