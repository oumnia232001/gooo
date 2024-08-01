package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/thedevsaddam/renderer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	rnd      *renderer.Render
)

// Informations d'identification pour la connexion à la base de données
func Init() {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Lire les variables d'environnement
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Construire la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Database = db
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw database connection: %v", err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalf("Failed to set Goose dialect: %v", err)
	}

	if err := goose.Up(dbSQL, "ressources/migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	// Initialize the renderer (  // Initialiser le moteur de rendu)
	rnd = renderer.New(renderer.Options{
		ParseGlobPattern: "./static/home.tpl",
	})
}

func GetRenderer() *renderer.Render {
	return rnd
}
