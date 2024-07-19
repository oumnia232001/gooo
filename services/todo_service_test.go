package services_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	return gormDB, mock, sqlDB
}

func TestDeleteTodoService(t *testing.T) {
	var sqlDB *sql.DB
	var mock sqlmock.Sqlmock
	var gormDB *gorm.DB

	testCases := []struct {
		name        string
		id          uint
		setup       func()
		checkResult func(err error)
	}{
		{
			name: "success",
			id:   1,
			setup: func() {
				gormDB, mock, sqlDB = initMockDB()
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			checkResult: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "database error",
			id:   1,
			setup: func() {
				gormDB, mock, sqlDB = initMockDB()
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			checkResult: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, gorm.ErrInvalidTransaction.Error(), err.Error())
			},
		},
		{
			name: "bad id",
			id:   0,
			setup: func() {
				// Pas de setup pour ce cas
			},
			checkResult: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "invalid ID", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			service := &services.TodoServiceImp{Db: gormDB}
			err := service.Delete(tc.id)
			tc.checkResult(err)
			if tc.setup != nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}

	defer sqlDB.Close()
}

func TestCreateTodoService(t *testing.T) {
	var sqlDB *sql.DB
	var mock sqlmock.Sqlmock
	var gormDB *gorm.DB

	testCases := []struct {
		name        string
		todo        models.TodoModel
		setup       func()
		checkResult func(todo models.TodoModel, err error)
	}{
		{
			name: "success",
			todo: models.TodoModel{
				Title:     "Test Todo",
				Completed: false,
			},
			setup: func() {
				gormDB, mock, sqlDB = initMockDB()
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `todo_models`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			checkResult: func(createdTodo models.TodoModel, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "Test Todo", createdTodo.Title)
				assert.False(t, createdTodo.Completed)
				assert.NotEqual(t, time.Time{}, createdTodo.CreatedAt)
			},
		},

		{
			name: "database error",
			todo: models.TodoModel{
				Title:     "Test Todo",
				Completed: false,
			},
			setup: func() {
				gormDB, mock, sqlDB = initMockDB()
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `todo_models`").WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			checkResult: func(todo models.TodoModel, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "todo not created", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			service := &services.TodoServiceImp{Db: gormDB}
			createdTodo, err := service.Create(tc.todo)
			tc.checkResult(createdTodo, err)
			if tc.setup != nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}

	defer sqlDB.Close()
}
