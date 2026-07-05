CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY , -- UUIDv4 strings for unguessable security IDs
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar_url TEXT,
    nickname TEXT UNIQUE,
    username TEXT UNIQUE,
    about_me TEXT,
    is_public INTEGER DEFAULT 1, -- 1 for public, 0 for private
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);