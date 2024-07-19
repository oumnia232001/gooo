package db_test

import (
	"testing"

	"github.com/go-todo1/db"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// la fonction Init pour initialiser la base de données
	db.Init()

	// Vérifier que Database est initialisé et non nul
	assert.NotNil(t, db.Database, "Expected Database to be initialized")
}

func TestGetRenderer(t *testing.T) {
	// Appeler la fonction Init pour initialiser la base de données
	db.Init()

	// Récupérer le renderer
	rnd := db.GetRenderer()

	// Vérifier que le renderer est initialisé et non nul
	assert.NotNil(t, rnd, "Expected Renderer to be initialized")
}

func TestGetDB(t *testing.T) {
	// la fonction Init pour initialiser la base de données
	db.Init()

	// Récupérer la base de données
	database := db.Database

	// Vérifier que la base de données est initialisée et non nulle
	assert.NotNil(t, database, "Expected Database to be initialized")
}
