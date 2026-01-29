package handlers

import (
	"database/sql"
	"net/http"
	"pastebin/internal/models"
	"strconv"
	"strings"
)

func LoadPasteByID(r *http.Request) (*models.Paste, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/view/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return models.GetPaste(db, id)
}

func LoadPasteByShort(r *http.Request) (*models.Paste, error) {
	short := strings.TrimPrefix(r.URL.Path, "/p/")
	if short == "" {
		return nil, sql.ErrNoRows
	}
	return models.GetPasteByShortLink(db, short)
}
