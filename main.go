package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	controllers "github.com/go-todo1/Controllers"
)

const port = ":9000"

func main() { // point d'entrée
	controllers.InitRenderAndDB() // Initialise le moteur de rendu et la base de données pour les contrôleurs
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)         // Ajoute un middleware au routeur
	r.Get("/", homeHandler)          // Enregistre la route de la page d'accueil
	r.Mount("/todo", todoHandlers()) // Sous-routeur pour les TODOs

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("server gracefully stopped!")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter() // Création d'un nouveau routeur

	// Définir les routes et les handlers correspondants
	rg.Group(func(r chi.Router) {
		r.Get("/", controllers.FetchTodos)
		r.Post("/", controllers.CreateTodo)
		r.Put("/{id}", controllers.UpdateTodo)
		r.Delete("/{id}", controllers.DeleteTodo)
	})
	return rg
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := controllers.GetRenderer().Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	if err != nil {
		log.Fatal(err)
	}
}
