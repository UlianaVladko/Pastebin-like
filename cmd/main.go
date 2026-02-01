package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"pastebin/internal/handlers"
	"pastebin/internal/middleware"
	"pastebin/internal/services"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// db, err := sql.Open("sqlite3", "./pastes.db")

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is empty")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// schema, _ := os.ReadFile("schema.sql")
	// db.Exec(string(schema))

	handlers.Init(db)
	services.StartExpiredPastesCleanup(db, 10*time.Second)

	http.HandleFunc("/", handlers.TmplPasteHandler)
	http.HandleFunc("/paste", handlers.RateLimit(handlers.CreatePasteHandler))
	http.HandleFunc(
		"/p/",
		middleware.PasteMiddleware(
			middleware.LoadPasteByShort(db),
			handlers.Render404,
			middleware.PasteAccessMiddleware(
				handlers.Render404,
				handlers.ViewPasteHandler,
			),
		),
	)
	http.HandleFunc("/edit/", handlers.EditPasteHandler)
	http.HandleFunc("/delete/", handlers.DeletePasteHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
