package db

import (
	x "myProject/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMySQLDB() (*gorm.DB, error) {
	// dsn := "user:password@tcp(mysql:3306)/Events?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:root@tcp(localhost:3306)/Events?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&x.Event{}); err != nil {
		return nil, err
	}

	return db, nil
}
