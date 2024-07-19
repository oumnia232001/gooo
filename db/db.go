package db

import (
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
	dsn := "root:12345678@tcp(127.0.0.1:3306)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Database = db

	// Initialize the renderer (  // Initialiser le moteur de rendu)
	rnd = renderer.New(renderer.Options{
		ParseGlobPattern: "./static/*.tpl",
	})
}

func GetRenderer() *renderer.Render {
	return rnd
}
