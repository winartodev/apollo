CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50) DEFAULT NULL,
    last_name VARCHAR(50) DEFAULT NULL,
    profile_picture VARCHAR(50) DEFAULT NULL,
    password VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) DEFAULT NULL,
    is_email_verified BOOL DEFAULT FALSE,
    is_phone_verified BOOL DEFAULT FALSE,
    last_login BIGINT DEFAULT 0,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);