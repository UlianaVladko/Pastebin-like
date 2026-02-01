package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"pastebin/internal/middleware"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/internal/utils"
)

var (
	db            *sql.DB
	baseTemplates = template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/header.html",
		),
	)
)

func renderTmpl(w http.ResponseWriter, contentFile string, data interface{}) {
	tmpl, err := baseTemplates.Clone()
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	_, err = tmpl.ParseFiles(contentFile)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func TmplPasteHandler(w http.ResponseWriter, r *http.Request) {
	token := generateCSRFToken()

	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: token,
		Path:  "/",
	})

	renderTmpl(w, "templates/form.html", map[string]string{
		"CSRFToken": token,
	})
}

func Init(database *sql.DB) {
	db = database
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

	// const maxPasteLength = 10_000
	content := r.FormValue("content")
	if len(content) == 0 {
		http.Error(w, "Paste cannot be empty", http.StatusBadRequest)
		return
	}
	// if len(content) > maxPasteLength {
	// 	http.Error(w, "Paste is too large", http.StatusBadRequest)
	// 	return
	// }

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "Untitled"
	}

	lang := r.FormValue("language")
	if lang == "" {
		lang = models.LangText
	}

	isPrivate := r.FormValue("private") == "true"

	var expiresAt *time.Time
	exp := r.FormValue("expires")
	if exp != "" && exp != "never" {
		expiresAt = services.ComputeExpiration(exp)
	}

	_, shortLink, editToken, err := models.CreatePaste(db, content, lang, name, isPrivate, expiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isPrivate {
		http.Redirect(w, r, "/p/"+shortLink+"?token="+editToken, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/p/"+shortLink, http.StatusSeeOther)
	}
}

func Render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	renderTmpl(w, "templates/err404.html", nil)
}

func ViewPasteHandler(w http.ResponseWriter, r *http.Request) {
	paste := r.Context().Value(middleware.PasteKey).(*models.Paste)
	editToken := r.URL.Query().Get("token")

	fileSize := len(paste.Content)
	humanSize := fmt.Sprintf("%.2f KB", float64(fileSize)/1024)

	renderTmpl(w, "templates/view.html", struct {
		*models.Paste
		CreatedAgo string
		EditToken  string
		Edited     bool
		FileSize   string
	}{
		Paste:      paste,
		CreatedAgo: utils.TimeAgo(paste.CreatedAt),
		EditToken:  editToken,
		Edited:     paste.UpdatedAt != nil && !paste.UpdatedAt.Equal(paste.CreatedAt),
		FileSize:   humanSize,
	})
}

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

	if r.Method == http.MethodGet {
		renderTmpl(w, "templates/edit.html", paste)
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

	http.Redirect(w, r, "/p/"+*paste.ShortLink+"?token="+token, http.StatusSeeOther)
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

	_, err = db.Exec(`DELETE FROM pastes WHERE id=$1`, paste.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
