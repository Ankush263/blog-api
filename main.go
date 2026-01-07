package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ankush263/blog-api/internal/db"
	"github.com/Ankush263/blog-api/internal/handler"
	"github.com/Ankush263/blog-api/internal/middleware"
	"github.com/Ankush263/blog-api/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	dsn := os.Getenv("CONNECTION_STRING")
	if dsn == "" {
		log.Fatal("CONNECTION_STRING not set")
	}

	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// repo := repository.NewPostRepository(dbConn)
	// handler := handler.NewPostHandler(repo)

	repo := repository.NewPostRepository(dbConn)
	handler := handler.NewPostHandler(repo)

	r := mux.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Recover)
	r.Use(middleware.CORS)

	r.HandleFunc("/posts", handler.Create).Methods("POST")
	r.HandleFunc("/posts", handler.GetAll).Methods("GET")
	r.HandleFunc("/posts/{id}", handler.GetById).Methods("GET")
	r.HandleFunc("/posts/{id}", handler.Update).Methods("PATCH")
	r.HandleFunc("/posts/{id}", handler.Delete).Methods("DELETE")

	log.Println("Server is running on PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
