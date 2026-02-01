CREATE TABLE IF NOT EXISTS pastes (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    name TEXT NOT NULL,
    language TEXT DEFAULT 'text',
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    short_link TEXT UNIQUE,
    edit_token TEXT UNIQUE,
    updated_at TIMESTAMP NOT NULL
);