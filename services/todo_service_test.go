package services_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initMockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return gormDB, mock, sqlDB, nil
}

func TestDeleteTodoService(t *testing.T) {
	testCases := []struct {
		name        string
		id          uint
		setup       func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error)
		checkResult func(err error)
	}{
		{
			name: "success",
			id:   1,
			setup: func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
				gormDB, mock, sqlDB, err := initMockDB()
				if err != nil {
					return nil, nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return gormDB, mock, sqlDB, nil
			},
			checkResult: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "database error",
			id:   1,
			setup: func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
				gormDB, mock, sqlDB, err := initMockDB()
				if err != nil {
					return nil, nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
				return gormDB, mock, sqlDB, nil
			},
			checkResult: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, gorm.ErrInvalidTransaction.Error(), err.Error())
			},
		},
		{
			name: "bad id",
			id:   0,
			setup: func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
				return nil, nil, nil, nil
			},
			checkResult: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "invalid ID", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gormDB, mock, sqlDB, err := tc.setup()
			if err != nil {
				t.Fatalf("Failed to set up test case: %v", err)
			}
			if gormDB == nil && tc.name != "bad id" {
				t.Fatalf("gormDB is nil")
			}
			service := services.NewTodoServiceImp(gormDB) // Utilisez NewTodoServiceImp ici
			err = service.Delete(tc.id)
			tc.checkResult(err)
			if mock != nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	}
}

func TestCreateTodoService(t *testing.T) {
	testCases := []struct {
		name        string
		todo        models.TodoModel
		setup       func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error)
		checkResult func(createdTodo models.TodoModel, err error)
	}{
		{
			name: "success",
			todo: models.TodoModel{
				Title:     "Test Todo",
				Completed: false,
			},
			setup: func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
				gormDB, mock, sqlDB, err := initMockDB()
				if err != nil {
					return nil, nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `todo_models`").
					WithArgs("Test Todo", false, sqlmock.AnyArg(), sqlmock.AnyArg()). // add arguments for timestamps
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return gormDB, mock, sqlDB, nil
			},
			checkResult: func(createdTodo models.TodoModel, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "Test Todo", createdTodo.Title)
				assert.Equal(t, false, createdTodo.Completed)
			},
		},
		{
			name: "database error",
			todo: models.TodoModel{
				Title:     "Test Todo",
				Completed: false,
			},
			setup: func() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
				gormDB, mock, sqlDB, err := initMockDB()
				if err != nil {
					return nil, nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `todo_models`").
					WithArgs("Test Todo", false, sqlmock.AnyArg(), sqlmock.AnyArg()). // add arguments for timestamps
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
				return gormDB, mock, sqlDB, nil
			},
			checkResult: func(createdTodo models.TodoModel, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, gorm.ErrInvalidTransaction.Error(), err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gormDB, mock, sqlDB, err := tc.setup()
			if err != nil {
				t.Fatalf("Failed to set up test case: %v", err)
			}
			if gormDB == nil {
				t.Fatalf("gormDB is nil")
			}
			service := services.NewTodoServiceImp(gormDB) // Utilisez NewTodoServiceImp ici
			createdTodo, err := service.Create(tc.todo)
			tc.checkResult(createdTodo, err)
			if mock != nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	}
}
