package db

import (
	x "myProject/admin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMySQLDB() (*gorm.DB, error) {
	dsn := "root:Venu@6582@tcp(127.0.0.1:3306)/Events?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&x.Event{}); err != nil {
		return nil, err
	}

	return db, nil
}
