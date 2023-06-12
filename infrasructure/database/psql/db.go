package database

import (
	"fmt"

	"gorm.io/driver/posgres"
	"gorm.io/gorm"
)

func NewPosgresConnection(conf config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.Username, conf.Name, conf.Port)
	db, err := gorm.Open(posgres.Open(dsn), gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ = db.DB()
	sqlDB.SetMaxIdleConn(conf.MaxIdleConnections)
	sqlDB.SetMaxOpenConn(conf.MaxOpenConnections)

}
