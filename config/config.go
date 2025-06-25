package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"perpustakaan-api/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost",    // DB_HOST
		"admin1",       // DB_USER
		"admin1234",   // DB_PASSWORD
		"perpustakaan_db", // DB_NAME
		"5432",        // DB_PORT
	)


	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.AdminDPK{},
		&models.AdminPerpustakaan{},
		&models.Perpustakaan{},
		&models.SDM{},
		&models.Pengunjung{},
		&models.Anggota{},
		&models.Verifikasi{},
		&models.Revisi{},
		&models.Laporan{},
		&models.LogAktivitas{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}