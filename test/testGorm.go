package main

import (
	"fmt"
	"gochat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/gochat?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})
	user := &models.UserBasic{}
	user.Name = "test"
	db.Create(user)

	// Read

	db.First(&user, 1)
	fmt.Println(db.First(user, 1))

	db.Model(user).Update("Password", "123456")
}
