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
	"github.com/go-todo1/db"
)

const port = ":9000"

func main() { //point d'entre
	db.Init()                     // Initialise la connexion à la base de données
	controllers.InitRenderAndDB() // Initialise le moteur de rendu et la base de données pour les contrôleurs
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)         //Ajoute un middleware au routeur.
	r.Get("/", homeHandler)          // enregistr la route de site
	r.Mount("/todo", todoHandlers()) //sous roteur sur la route

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
	rg := chi.NewRouter() // new router creation

	//  et ici Définir les routes et les handlers correspondants
	rg.Group(func(r chi.Router) {
		r.Get("/", controllers.FetchTodos)
		r.Post("/", controllers.CreateTodo)
		r.Put("/{id}", controllers.UpdateTodo)
		r.Delete("/{id}", controllers.DeleteTodo)
	})
	return rg
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rnd := db.GetRenderer() // récupération  de moteur
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	if err != nil {
		log.Fatal(err)
	}
} //fonction qu'il répond aux requêtes HTTP.
