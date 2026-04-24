-- Active: 1774007783901@@127.0.0.1@5432@orderhub
CREATE Table IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
)

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);