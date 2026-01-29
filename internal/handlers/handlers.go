package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pastebin/internal/middleware"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/internal/utils"
)

var (
	db        *sql.DB
	templates = template.Must(template.ParseGlob("templates/*.html"))
)

func Init(database *sql.DB) {
	db = database
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	token := generateCSRFToken()

	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: token,
		Path:  "/",
	})

	templates.ExecuteTemplate(w, "form.html", token)
}

func CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	formToken := r.FormValue("csrf_token")
	cookie, err := r.Cookie("csrf_token")
	if err != nil || cookie.Value != formToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	const maxPasteLength = 10_000
	content := r.FormValue("content")
	if len(content) == 0 {
		http.Error(w, "Paste cannot be empty", http.StatusBadRequest)
		return
	}
	if len(content) > maxPasteLength {
		http.Error(w, "Paste is too large", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "Untitled"
	}

	lang := r.FormValue("language")
	if lang == "" {
		lang = models.LangText
	}

	isPrivate := r.FormValue("private") == "true"

	exp := r.FormValue("expires")
	var expiresAt *time.Time
	if exp != "" && exp != "never" {
		expiresAt = services.ComputeExpiration(exp)
	}

	id, editToken, err := models.CreatePaste(db, content, lang, name, isPrivate, expiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortLink := strconv.FormatInt(id, 10)
	editURL := "/edit/" + shortLink + "?token=" + editToken
	log.Println("EDIT URL:", editURL)
	log.Println("EDIT TOKEN:", editToken)

	http.Redirect(w, r, "/view/"+strconv.FormatInt(id, 10)+"?edit_token="+editToken, http.StatusSeeOther)
}

func Render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "err404.html", nil)
}

func ViewPasteHandler(w http.ResponseWriter, r *http.Request) {

	paste := r.Context().Value(middleware.PasteKey).(*models.Paste)

	editToken := r.URL.Query().Get("edit_token")

	templates.ExecuteTemplate(w, "view.html", struct {
		*models.Paste
		CreatedAgo string
		EditToken  string
		Edited     bool
	}{
		Paste:      paste,
		CreatedAgo: utils.TimeAgo(paste.CreatedAt),
		EditToken:  editToken,
		Edited:     paste.UpdatedAt != nil && !paste.UpdatedAt.Equal(paste.CreatedAt),
	})
}

// func ViewPasteShortHandler(w http.ResponseWriter, r *http.Request) {
// 	short := strings.TrimPrefix(r.URL.Path, "/p/")
// 	if short == "" {
// 		render404(w)
// 		return
// 	}

// 	paste, err := models.GetPasteByShortLink(db, short)
// 	if err != nil || paste == nil {
// 		render404(w)
// 		return
// 	}

// 	if paste.ExpiresAt != nil && time.Now().UTC().After(*paste.ExpiresAt) {
// 		render404(w)
// 		return
// 	}

// 	templates.ExecuteTemplate(w, "view.html", struct {
// 		*models.Paste
// 		CreatedAgo string
// 		Edited     bool
// 	}{
// 		Paste:      paste,
// 		CreatedAgo: utils.TimeAgo(paste.CreatedAt),
// 		Edited:     paste.UpdatedAt != nil && !paste.UpdatedAt.Equal(paste.CreatedAt),
// 	})
// }

func EditPasteHandler(w http.ResponseWriter, r *http.Request) {
	short := strings.TrimPrefix(r.URL.Path, "/edit/")
	token := r.URL.Query().Get("token")
	if short == "" || token == "" {
		Render404(w)
		return
	}

	paste, err := models.GetPasteByShortLink(db, short)
	if err != nil || paste == nil || paste.EditToken != token {
		Render404(w)
		return
	}

	if paste.ExpiresAt != nil && time.Now().UTC().After(*paste.ExpiresAt) {
		Render404(w)
		return
	}

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "edit.html", paste)
		return
	}

	content := r.FormValue("content")
	name := r.FormValue("name")
	lang := r.FormValue("language")
	isPrivate := r.FormValue("private") == "true"

	err = models.UpdatePaste(db, paste.ID, content, name, lang, isPrivate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+strconv.FormatInt(paste.ID, 10), http.StatusSeeOther)
}

func DeletePasteHandler(w http.ResponseWriter, r *http.Request) {
	short := strings.TrimPrefix(r.URL.Path, "/delete/")
	token := r.URL.Query().Get("token")
	if short == "" || token == "" {
		Render404(w)
		return
	}

	paste, err := models.GetPasteByShortLink(db, short)
	if err != nil || paste == nil || paste.EditToken != token {
		Render404(w)
		return
	}

	_, err = db.Exec(`DELETE FROM pastes WHERE id=?`, paste.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
