package services

import (
	"errors"
	"time"

	models "github.com/go-todo1/Models"
	"github.com/go-todo1/db"
	"gorm.io/gorm"
)

type TodoService interface {
	Create(todo models.TodoModel) (models.TodoModel, error)
	Update(id uint, todo models.TodoModel) (models.TodoModel, error)
	Delete(id uint) error
}

type todoServiceImp struct {
	db *gorm.DB
}

func NewTodoService() TodoService {
	return &todoServiceImp{
		db: db.GetDB(),
	}
}

func (tds *todoServiceImp) Create(todo models.TodoModel) (models.TodoModel, error) {
	todo.CreatedAt = time.Now()
	if err := tds.db.Create(&todo).Error; err != nil {
		return models.TodoModel{}, err

	}
	return todo, nil
}

func (tds *todoServiceImp) Update(id uint, todo models.TodoModel) (models.TodoModel, error) {
	var existingTodo models.TodoModel
	if err := tds.db.First(&existingTodo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoModel{}, errors.New("todo not found")
		}
		return models.TodoModel{}, err
	}

	todo.UpdatedAt = time.Now()
	if err := tds.db.Model(&existingTodo).Updates(todo).Error; err != nil {
		return models.TodoModel{}, err
	}
	return existingTodo, nil
}

func (tds *todoServiceImp) Delete(id uint) error {
	if err := tds.db.Delete(&models.TodoModel{}, id).Error; err != nil {
		return err
	}
	return nil
}
