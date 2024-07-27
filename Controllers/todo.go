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
var todoService services.TodoService
var Database *gorm.DB

func InitRenderAndDB() {
	InitDatabase()
	rnd = renderer.New()
	todoService = services.NewTodoServiceImp(Database)
}

func InitDatabase() {
	var err error
	dsn := "root:197520012003@tcp(mysql:3306)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func GetRenderer() *renderer.Render {
	return rnd
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
