package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	rnd      *renderer.Render
)

type Config struct {
	DBUser string `yaml:"DB_USER"`
	DBHost string `yaml:"DB_HOST"`
	DBPort string `yaml:"DB_PORT"`
	DBName string `yaml:"DB_NAME"`
}

func Init() {
	//  les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	//  les variables d'environnement pour le mot de passe
	password := os.Getenv("DB_PASSWORD")

	// Charger la configuration depuis le fichier config.yaml
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error opening config.yaml file: %v", err)
	}
	defer file.Close()

	config := Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config.yaml file: %v", err)
	}

	// Construire la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, password, config.DBHost, config.DBPort, config.DBName)
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

	// Initialize the renderer (Initialiser le moteur de rendu)
	rnd = renderer.New(renderer.Options{
		ParseGlobPattern: "./static/home.tpl",
	})
}

func GetRenderer() *renderer.Render {
	return rnd
}
