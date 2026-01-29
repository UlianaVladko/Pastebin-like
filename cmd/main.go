package main

import (
	"database/sql"
	"log"
	"net/http"
	"pastebin/internal/handlers"
	"pastebin/internal/middleware"
	"pastebin/internal/services"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./pastes.db")
	if err != nil {
		log.Fatal(err)
	}

	// schema, _ := os.ReadFile("schema.sql")
	// db.Exec(string(schema))

	handlers.Init(db)

	services.StartExpiredPastesCleanup(db, 10*time.Second)

	http.HandleFunc("/", handlers.FormHandler)
	http.HandleFunc("/paste", handlers.RateLimit(handlers.CreatePasteHandler))
	http.HandleFunc(
		"/view/",
		middleware.PasteMiddleware(handlers.LoadPasteByID, handlers.Render404, handlers.ViewPasteHandler),
	)
	http.HandleFunc(
		"/p/",
		middleware.PasteMiddleware(handlers.LoadPasteByShort, handlers.Render404, handlers.ViewPasteHandler),
	)
	http.HandleFunc("/edit/", handlers.EditPasteHandler)
	http.HandleFunc("/delete/", handlers.DeletePasteHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
