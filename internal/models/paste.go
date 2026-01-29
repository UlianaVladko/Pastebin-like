package models

import (
	"database/sql"
	"math/rand"
	"time"
)

const (
	LangText = "text"
	LangGo   = "go"
	LangJS   = "js"
	LangSQL  = "sql"
)

type Paste struct {
	ID        int64
	NamePaste string
	Content   string
	Language  string
	CreatedAt time.Time
	ExpiresAt *time.Time
	UpdatedAt *time.Time
	IsPrivate bool
	ShortLink *string
	EditToken string
}

func CreatePaste(db *sql.DB, content, language, name string, isPrivate bool, expiresAt *time.Time) (int64, string, error) {
	shortLink := generateShortCode(8)
	editToken := generateShortToken()
	now := time.Now().UTC()

	result, err := db.Exec(
		`INSERT INTO pastes (content, language, name, is_private, expires_at, created_at, updated_at, short_link, edit_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		content, language, name, boolToInt(isPrivate), expiresAt, now, now, shortLink, editToken,
	)
	if err != nil {
		return 0, "", err
	}
	id, err := result.LastInsertId()
	return id, editToken, err
}

func boolToInt(isPrivate bool) int {
	if isPrivate {
		return 1
	}
	return 0
}

func GetPaste(db *sql.DB, id int64) (*Paste, error) {
	row := db.QueryRow(
		`SELECT id, name, content, created_at, updated_at, is_private, expires_at, language, short_link, edit_token FROM pastes WHERE id = ?`,
		id,
	)

	var p Paste
	var isPrivateInt int
	err := row.Scan(
		&p.ID,
		&p.NamePaste,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
		&isPrivateInt,
		&p.ExpiresAt,
		&p.Language,
		&p.ShortLink,
		&p.EditToken,
	)
	if err != nil {
		return nil, err
	}
	p.IsPrivate = isPrivateInt != 0
	return &p, nil
}

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = base62[rand.Intn(len(base62))]
	}
	return string(b)
}

func generateShortToken() string {
	return generateShortCode(32)
}

func GetPasteByShortLink(db *sql.DB, short string) (*Paste, error) {
	row := db.QueryRow(
		`SELECT id, name, content, created_at, updated_at, is_private, expires_at, language, short_link, edit_token FROM pastes WHERE short_link = ?`,
		short,
	)

	var p Paste
	var isPrivateInt int
	err := row.Scan(
		&p.ID,
		&p.NamePaste,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
		&isPrivateInt,
		&p.ExpiresAt,
		&p.Language,
		&p.ShortLink,
		&p.EditToken,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.IsPrivate = isPrivateInt != 0
	return &p, nil
}

func UpdatePaste(db *sql.DB, id int64, content, name, language string, isPrivate bool) error {
	now := time.Now().UTC()
	_, err := db.Exec(
		`UPDATE pastes SET content=?, name=?, language=?, is_private=?, updated_at=? WHERE id=?`,
		content, name, language, boolToInt(isPrivate), now, id,
	)
	return err
}

func DeleteExpiredPastes(db *sql.DB) error {
	_, err := db.Exec(
		`DELETE FROM pastes WHERE expires_at IS NOT NULL AND expires_at <= CURRENT_TIMESTAMP`,
	)
	return err
}
