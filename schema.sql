CREATE TABLE IF NOT EXISTS pastes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    name TEXT NOT NULL,
    language TEXT DEFAULT 'text',
    is_private INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
    short_link TEXT UNIQUE,
    edit_token TEXT UNIQUE,
    updated_at DATETIME NOT NULL
);