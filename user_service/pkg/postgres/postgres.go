package postgres

import (
	"Booking_system/user_service/internal/config"
	"Booking_system/user_service/internal/models"
	"Booking_system/user_service/internal/util"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(`host=%s dbname=%s user=%s password=%s port=%d sslmode=disable`, config.Cfg.DataBase.Host, config.Cfg.DataBase.DBName, config.Cfg.DataBase.User, config.Cfg.DataBase.Password, config.Cfg.DataBase.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connected Postgres")
	}
	db.AutoMigrate(&models.User{})
	return db
}

func Paginate(pagination util.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
