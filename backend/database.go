package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal menghubungkan database:", err)
	}

	DB = database

	err = DB.AutoMigrate(
    &User{},
    &Order{},
    &OrderEvent{},
    &CancellationRequest{},
)

if err != nil {
	log.Fatal("Gagal membuat tabel:", err)
}

fmt.Println("Tabel database berhasil dibuat!")

	fmt.Println("PostgreSQL berhasil terhubung!")
}