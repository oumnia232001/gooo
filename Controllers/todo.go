package Controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/renderer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var rnd *renderer.Render
var todoService services.TodoService
var Database *gorm.DB

func InitRenderAndDB() {
	InitDatabase()
	rnd = renderer.New(renderer.Options{
		ParseGlobPattern: "static/*.tpl",
	})
	// Charger la clé API depuis les variables d'environnement
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		log.Fatal("RAPIDAPI_KEY environment variable not set")
	}
	todoService = services.NewTodoServiceImp(Database, apiKey)
}

func InitDatabase() {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Lire les variables de mot de passe
	password := os.Getenv("DB_PASSWORD")

	// Configurer viper pour lire le fichier config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Lire les variables de configuration
	user := viper.GetString("DB_USER")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	name := viper.GetString("DB_NAME")

	// Construire la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	Database = db

	// Assurer la migration de la base de données
	if !Database.Migrator().HasTable(&models.TodoModel{}) {
		if err := Database.Migrator().CreateTable(&models.TodoModel{}); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}

	if err := Database.AutoMigrate(&models.TodoModel{}); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database connected and tables ensured")
}

func GetRenderer() *renderer.Render {
	return rnd
}

type MockRenderer struct{}

func (m *MockRenderer) JSON(w http.ResponseWriter, status int, v interface{}, headers ...http.Header) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte("mock json response"))
	return err
}

func (m *MockRenderer) Template(w http.ResponseWriter, status int, files []string, data interface{}, headers ...http.Header) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte("mock template"))
	return err
}

func FetchTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.TodoModel
	if err := Database.Find(&todos).Error; err != nil {
		log.Printf("Error fetching todos: %v", err)
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err.Error(),
		})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todos,
	})
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			log.Printf("Error decoding JSON: %v", err)
			rnd.JSON(w, http.StatusBadRequest, renderer.M{
				"message": "Invalid request payload",
				"error":   err.Error(),
			})
			return
		}
		log.Printf("Error decoding JSON: %v", err)
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}
	if t.ID != 0 {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "ID should not be provided",
		})
		return
	}
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})
		return
	}
	createdTodo, err := todoService.Create(t)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to save todo",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo":    createdTodo,
	})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := todoService.Delete(uint(id)); err != nil {
		log.Printf("Error deleting todo: %v", err)
		if err.Error() == "database error" {
			http.Error(w, "database error", http.StatusInternalServerError)
		} else {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo deleted successfully"))
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid ID: %v", idParam)
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid ID",
		})
		return
	}

	var t models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Updating Todo with ID: %d and Data: %+v", id, t)

	updatedTodo, err := todoService.Update(uint(id), t)
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to update todo",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
		"todo":    updatedTodo,
	})
}

// Nouvelle fonction pour obtenir une citation depuis l'API RapidAPI
func GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
	// Appeler la méthode GetQuote du service
	quote, err := todoService.GetQuote()
	if err != nil {
		log.Printf("Error fetching quote: %v", err)
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to fetch quote",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusOK, quote)
}
