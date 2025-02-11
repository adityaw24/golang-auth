package databases

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGormDB(db *sql.DB) *gorm.DB {

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	fmt.Println("Connected to the database!")

	return gormDB
}
