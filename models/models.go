package models

import (
	"bunker-web/configs"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func init() {
	fmt.Println("Connecting to database...")
	connect()

	fmt.Println("Database connected, migrating tables...")
	// Create tables if not exists
	DB.AutoMigrate(
		&Announcement{},
		&Log{},
		&AndroidMpayUser{},
		&WindowsMpayUser{},
		&UnlimitedRentalServer{},
		&UserBanRecord{},
		&User{},
		&WebAuthnCredential{},
	)

	fmt.Println("Database initialised")
}

func connect() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.DB_USER, configs.DB_PASSWORD, configs.DB_HOST, configs.DB_NAME,
	)

	maxRetries := 5

	for i := range maxRetries {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(configs.GORM_LOGGER_MODE),
		})
		if err == nil {
			DB = db
			break
		}
		if i == maxRetries-1 {
			panic(err)
		}
		fmt.Println("Database connection failed, retrying...")
		time.Sleep(3 * time.Second)
	}

	// Set connection pool
	d, _ := DB.DB()
	d.SetMaxIdleConns(10)
	d.SetMaxOpenConns(100)
	d.SetConnMaxIdleTime(time.Minute)
	d.SetConnMaxLifetime(time.Hour)

	// Heartbeat
	go heartbeat(d)
}

// Heartbeat to keep connection alive
func heartbeat(d *sql.DB) {
	for {
		if d.Ping() != nil {
			fmt.Println("Database connection lost, reconnecting...")
			connect()
			break
		}
		time.Sleep(30 * time.Second)
	}
}

func DBCreate(entity any) error {
	return DB.Create(entity).Error
}

func DBSave(entity any) error {
	return DB.Save(entity).Error
}

func DBDelete(entity any) error {
	return DB.Delete(entity).Error
}

func DBRemove(entity any) error {
	return DB.Unscoped().Delete(entity).Error
}
