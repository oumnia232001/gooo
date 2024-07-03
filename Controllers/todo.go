package Controllers

import (
	"encoding/json" //convertir les donnees go en json et vice versa
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/db"
	"github.com/thedevsaddam/renderer"
	"gorm.io/gorm" // Il permet de manipuler la base de données  on utilisant go au lieu de sql
)

var rnd *renderer.Render
var Database *gorm.DB

func InitRenderAndDB() {
	rnd = db.GetRenderer()
	Database = db.GetDB()
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
		log.Printf("Error decoding JSON: %v", err) // Log the error
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})
		return
	}
	t.CreatedAt = time.Now()
	if err := Database.Create(&t).Error; err != nil {
		log.Printf("Error creating todo: %v", err) // Log the error
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}
	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "todo created successfully",
		"todo_id": t.ID,
	})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Printf("error : %s", err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := Database.Delete(&models.TodoModel{}, id).Error; err != nil {
		fmt.Printf("error : %s", err.Error())
		http.Error(w, "Error deleting todo with ID "+idStr+": "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo deleted successfully"))
	return
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var t models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := Database.Model(&models.TodoModel{}).Where("id = ?", id).Updates(t).Error; err != nil {
		log.Printf("Error updating todo with ID %s: %v", id, err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}
	log.Printf("Todo with ID %s updated successfully", id)
	response := map[string]string{"message": "todo updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
