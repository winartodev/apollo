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

-- Data Seeding
DO $$
    DECLARE current_epoch_time BIGINT;
    BEGIN
        current_epoch_time =  EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);

        INSERT INTO users (uuid, email, phone_number, username, profile_picture, first_name, last_name, password, is_email_verified, is_phone_verified, created_at, updated_at)
        VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'public-user@gmail.com', '0812345678', 'Public User', '', 'Public', 'User', '$2a$10$JLqd1c/C5Aj13v88fcMx8.tcue8JiOh9sEL8FRhtLy.l2Vn3.VBWm', true, true, current_epoch_time, current_epoch_time),
               ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'internal-user@gmail.com', '0823456789', 'Internal User', '', 'Internal', 'User', '$2a$10$JLqd1c/C5Aj13v88fcMx8.tcue8JiOh9sEL8FRhtLy.l2Vn3.VBWm', true, true, current_epoch_time, current_epoch_time),
               ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'protected-user@gmail.com', '0834567890', 'Protected User', '', 'Protected', 'User', '$2a$10$JLqd1c/C5Aj13v88fcMx8.tcue8JiOh9sEL8FRhtLy.l2Vn3.VBWm', true, true, current_epoch_time, current_epoch_time);
END$$