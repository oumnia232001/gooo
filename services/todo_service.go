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

type TodoServiceImp struct {
	Db *gorm.DB
}

func NewTodoService() TodoService {
	return &TodoServiceImp{
		Db: db.GetDB(),
	}
}

func (tds *TodoServiceImp) Create(todo models.TodoModel) (models.TodoModel, error) {
	todo.CreatedAt = time.Now()
	tx := tds.Db.Begin()
	if tx.Error != nil {
		// !=different  de (Opérateur de comparaison)
		return models.TodoModel{}, tx.Error
	}
	if err := tx.Create(&todo).Error; err != nil {
		// ; fin de condition
		tx.Rollback()
		//Rollback pout annuler transaction
		return models.TodoModel{}, errors.New("todo not created")
	}
	return todo, tx.Commit().Error
}

//tx : transaction

func (tds *TodoServiceImp) Update(id uint, todo models.TodoModel) (models.TodoModel, error) {
	var existingTodo models.TodoModel
	if err := tds.Db.First(&existingTodo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoModel{}, errors.New("todo not found")
		}
		return models.TodoModel{}, err
	}

	todo.UpdatedAt = time.Now()
	if err := tds.Db.Model(&existingTodo).Updates(todo).Error; err != nil {
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
