package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	models "github.com/go-todo1/Models"
	"gorm.io/gorm"
)

type TodoService interface {
	Create(todo models.TodoModel) (models.TodoModel, error)
	Update(id uint, todo models.TodoModel) (models.TodoModel, error)
	Delete(id uint) error
	GetQuote() (models.QuoteResponse, error)
}

func NewTodoServiceImp(db *gorm.DB, apiKey string) *TodoServiceImp {
	return &TodoServiceImp{Db: db, APIKey: apiKey}
}

type TodoServiceImp struct {
	Db      *gorm.DB
	APIKey  string
	Service TodoService
}

// Create
func (s *TodoServiceImp) Create(todo models.TodoModel) (models.TodoModel, error) {
	if todo.ID != 0 {
		return models.TodoModel{}, fmt.Errorf("invalid ID")
	}

	fmt.Println(&todo)
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&todo).Error; err != nil {
			return err
		}
		return nil
	})

	return todo, err
}

func (s *TodoServiceImp) Update(id uint, todo models.TodoModel) (models.TodoModel, error) {
	var existingTodo models.TodoModel
	if err := s.Db.First(&existingTodo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoModel{}, errors.New("todo not found")
		}
		return models.TodoModel{}, err
	}

	todo.UpdatedAt = time.Now()
	if err := s.Db.Model(&existingTodo).Updates(todo).Error; err != nil {
		return models.TodoModel{}, err
	}
	return existingTodo, nil
}

func (s *TodoServiceImp) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	tx := s.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Delete(&models.TodoModel{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Obtient une citation aléatoire de l'API
func (s *TodoServiceImp) GetQuote() (models.QuoteResponse, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("X-RapidAPI-Key", s.APIKey).
		SetHeader("X-RapidAPI-Host", "quotes15.p.rapidapi.com").
		Get("https://quotes15.p.rapidapi.com/quotes/random/?language_code=en")

	if err != nil {
		return models.QuoteResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return models.QuoteResponse{}, fmt.Errorf("unexpected response code %d", resp.StatusCode())
	}

	var quoteResp models.QuoteResponse
	if err := json.Unmarshal(resp.Body(), &quoteResp); err != nil {
		return models.QuoteResponse{}, err
	}

	return quoteResp, nil
}
