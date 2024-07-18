package Controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/services"
	"github.com/thedevsaddam/renderer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var rnd *renderer.Render = renderer.New()
var todoService services.TodoService = services.NewTodoService()
var Database *gorm.DB

func InitDatabase() {
	var err error
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func FetchTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.TodoModel
	if err := Database.Find(&todos).Error; err != nil {
		log.Printf("Error fetching todos: %v", err) // Log the error
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
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
				"message": "Invalid ID",
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
		"message": "todo created successfully",
		"todo_id": createdTodo.ID,
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
		http.Error(w, "Error deleting todo with ID "+idStr+": "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo deleted successfully"))
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var t models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedTodo, err := todoService.Update(uint(id), t)
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}
	log.Printf("Todo with ID %d updated successfully", id)
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "todo updated successfully",
		"todo":    updatedTodo,
	})
}
