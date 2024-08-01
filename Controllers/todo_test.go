package Controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	models "github.com/go-todo1/Models"
	"github.com/go-todo1/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thedevsaddam/renderer"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMockDB() *gorm.DB {
	dialector := sqlite.Open("file::memory:?cache=shared&_fk=1")
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&models.TodoModel{})
	db.Create(&models.TodoModel{Title: "Learn Go"})
	return db
}

// 		t.Errorf("unexpected response message: got %v want %v", response["message"], "todo created successfully")
// 	}
// }

func TestCreateTodo(t *testing.T) {
	rnd = renderer.New(renderer.Options{})
	defaultBody := `{
        "title": "sara"
    }`
	databaseErrorBody := `{
        "title": "sara"
    }`
	badIDBody := `{
        "title": "sara",
        "id": "invalid"
    }`

	testCases := []struct {
		name        string
		request     func() *http.Request
		setup       func()
		checkResult func(code int, response string)
	}{
		{
			name: "success",
			request: func() *http.Request {
				req, _ := http.NewRequest("POST", "/todo", strings.NewReader(defaultBody))
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Create(gomock.Eq(models.TodoModel{
					Title: "sara",
				})).Return(models.TodoModel{
					ID:        1,
					Title:     "sara",
					Completed: false,
				}, nil)
				todoService = todoServiceMock
			},
			checkResult: func(code int, response string) {
				assert.Equal(t, http.StatusCreated, code)
			},
		},
		{
			name: "database error",
			request: func() *http.Request {
				req, _ := http.NewRequest("POST", "/todo", strings.NewReader(databaseErrorBody))
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Create(gomock.Eq(models.TodoModel{
					Title: "sara",
				})).Return(models.TodoModel{}, errors.New("database error"))
				todoService = todoServiceMock
			},
			checkResult: func(code int, response string) {
				assert.Equal(t, http.StatusInternalServerError, code)
				assert.Contains(t, response, "database error")
			},
		},
		{
			name: "bad id",
			request: func() *http.Request {
				req, _ := http.NewRequest("POST", "/todo", strings.NewReader(badIDBody))
				return req
			},
			setup: func() {
				// No need to configure a mock because the invalid ID must be processed before calling the service
			},
			checkResult: func(code int, response string) {
				assert.Equal(t, http.StatusBadRequest, code)
				assert.Contains(t, response, "Invalid request payload")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			rr := httptest.NewRecorder()
			CreateTodo(rr, tc.request())
			tc.checkResult(rr.Code, rr.Body.String())
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	var sqlDB *sql.DB

	testCases := []struct {
		name        string
		id          string
		request     func() *http.Request
		setup       func()
		checkResult func(rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			id:   "1",
			request: func() *http.Request {
				req, _ := http.NewRequest("DELETE", "/todos/1", nil)
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Delete(gomock.Eq(uint(1))).Return(nil)
				todoService = todoServiceMock
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, "Todo deleted successfully", rr.Body.String())
			},
		},
		{
			name: "database error",
			id:   "1",
			request: func() *http.Request {
				req, _ := http.NewRequest("DELETE", "/todos/1", nil)
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Delete(gomock.Eq(uint(1))).Return(errors.New("database error"))
				todoService = todoServiceMock
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			rr := httptest.NewRecorder()
			router := chi.NewRouter()
			router.Delete("/todos/{id}", DeleteTodo)
			req := tc.request()
			router.ServeHTTP(rr, req)
			tc.checkResult(rr)
		})
	}

	if sqlDB != nil {
		defer sqlDB.Close()
	}
}

func TestUpdateTodo(t *testing.T) {
	rnd = renderer.New(renderer.Options{})
	testCases := []struct {
		name        string
		id          string
		request     func() *http.Request
		setup       func()
		checkResult func(rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			id:   "1",
			request: func() *http.Request {
				body := `{"title":"Updated Title","completed":true}`
				req, _ := http.NewRequest("PUT", "/todos/1", strings.NewReader(body))
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Update(uint(1), models.TodoModel{
					Title:     "Updated Title",
					Completed: true,
				}).Return(models.TodoModel{
					ID:        1,
					Title:     "Updated Title",
					Completed: true,
				}, nil)
				todoService = todoServiceMock
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Contains(t, rr.Body.String(), "Todo updated successfully")
			},
		},
		{
			name: "database error",
			id:   "1",
			request: func() *http.Request {
				body := `{"title":"Updated Title","completed":true}`
				req, _ := http.NewRequest("PUT", "/todos/1", strings.NewReader(body))
				return req
			},
			setup: func() {
				ctrl := gomock.NewController(t)
				todoServiceMock := mocks.NewMockTodoService(ctrl)
				todoServiceMock.EXPECT().Update(gomock.Eq(uint(1)), gomock.Eq(models.TodoModel{
					Title:     "Updated Title",
					Completed: true,
				})).Return(models.TodoModel{}, fmt.Errorf("database error"))
				todoService = todoServiceMock
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "Failed to update todo")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			rr := httptest.NewRecorder()
			router := chi.NewRouter()
			router.Put("/todos/{id}", UpdateTodo)
			req := tc.request()
			router.ServeHTTP(rr, req)
			tc.checkResult(rr)
		})
	}
}
