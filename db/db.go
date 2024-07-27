package db

import (
	"log"

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
	dsn := "root:197520012003@tcp(mysql:3306)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
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
