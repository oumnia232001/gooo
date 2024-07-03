package Controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	models "github.com/go-todo1/Models"
	"github.com/stretchr/testify/assert"
	"github.com/thedevsaddam/renderer"
	"gorm.io/driver/mysql"
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

func TestCreateTodo(t *testing.T) {
	// Initialisation des dépendances
	sqlDB, mock := initMockDB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `todo_models` \\(`title`,`completed`,`created_at`,`updated_at`\\) VALUES \\(\\?,\\?,\\?,\\?\\)$").WithArgs(
		"Learn Go",
		false,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var jsonStr = []byte(`{"Title":"Learn Go"}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTodo)

	// Initialiser `rnd` si nécessaire
	if rnd == nil {
		rnd = renderer.New()
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["message"] != "todo created successfully" {
		t.Errorf("unexpected response message: got %v want %v", response["message"], "todo created successfully")
	}
}

func TestDeleteTodo(t *testing.T) {
	var sqlDB *sql.DB
	var mock sqlmock.Sqlmock

	testCases := []struct {
		name        string
		id          string
		setup       func()
		checkResult func(rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			id:   "1",
			setup: func() {
				// Initialisation des dépendances
				sqlDB, mock = initMockDB()

				mock.ExpectBegin()
				// Passer un int64 comme argument attendu
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, "Todo deleted successfully", rr.Body.String())
			},
		},
		{
			name:  "database error",
			id:    "1",
			setup: func() {},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "Error deleting todo with ID")
			},
		},
		{
			name: "bad id",
			id:   "oumnia",
			setup: func() {
				// Initialisation des dépendances
				sqlDB, mock = initMockDB()

				mock.ExpectBegin()
				// Passer un int64 comme argument attendu
				mock.ExpectExec("^DELETE FROM `todo_models` WHERE `todo_models`.`id` = \\?$").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			checkResult: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Contains(t, rr.Body.String(), "Invalid ID")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			rr := httptest.NewRecorder()
			router := chi.NewRouter()
			router.Delete("/todos/{id}", DeleteTodo)
			req, _ := http.NewRequest("DELETE", "/todos/"+tc.id, nil)
			router.ServeHTTP(rr, req)
			tc.checkResult(rr)
		})
	}

	defer sqlDB.Close()
}
func TestUpdateTodo(t *testing.T) {
	// Initialisation des dépendances
	sqlDB, mock := initMockDB()
	defer sqlDB.Close()

	// Définition les attentes pour la base de données simulée
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `todo_models` SET `title`=?, `completed`=?, `created_at`=?, `updated_at`=? WHERE `id` = \\?$").
		WithArgs("Learn Go Updated", false, sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Créeation d'une requête HTTP PUT
	var jsonStr = []byte(`{"Title":"Learn Go Updated","Completed":false}`)
	req, err := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/todos/{id}", UpdateTodo)

	// Initialisation `rnd` si nécessaire
	if rnd == nil {
		rnd = renderer.New()
	}
	router.ServeHTTP(rr, req)

	// Vérification le code de réponse
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Vérification le contenu de la réponse
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	expectedMessage := "todo updated successfully"
	if response["message"] != expectedMessage {
		t.Errorf("unexpected response message: got %v want %v", response["message"], expectedMessage)
	}
}

func initMockDB() (*sql.DB, sqlmock.Sqlmock) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		DriverName:                "mysql",
		Conn:                      mdb,
	})

	Database, _ = gorm.Open(dialector, &gorm.Config{})
	return mdb, mock
}
