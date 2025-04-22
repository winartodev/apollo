-- Tables
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    is_active BOOL DEFAULT TRUE,
    description TEXT DEFAULT  '',
    created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_services_is_active ON services(is_active);
CREATE INDEX IF NOT EXISTS idx_services_created_at ON services(created_at);
CREATE INDEX IF NOT EXISTS idx_services_updated_at ON services(updated_at);
CREATE INDEX IF NOT EXISTS idx_services_active_created ON services(is_active, created_at);
CREATE INDEX IF NOT EXISTS idx_services_created_by ON services(created_by);
CREATE INDEX IF NOT EXISTS idx_services_updated_by ON services(updated_by);

-- Seeding Data
DO $$
    DECLARE current_epoch_time BIGINT;
    DECLARE super_user_id INT;
BEGIN
    current_epoch_time =  EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);

    INSERT INTO users (uuid, email, phone_number, username, profile_picture, first_name, last_name, password, is_email_verified, is_phone_verified, created_at, updated_at)
        VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a999', 'super-user@gmail.com', '08999999999', 'Super User', '', 'Super', 'User', '$2a$10$JLqd1c/C5Aj13v88fcMx8.tcue8JiOh9sEL8FRhtLy.l2Vn3.VBWm', true, true, current_epoch_time, current_epoch_time)
        RETURNING id INTO super_user_id;

    INSERT INTO services (slug, name, is_active, created_by, updated_by, created_at, updated_at)
        VALUES ('service-a', 'Services A', true,super_user_id, super_user_id, current_epoch_time, current_epoch_time),
               ('service-b', 'Services B', true,super_user_id, super_user_id, current_epoch_time, current_epoch_time),
               ('service-c', 'Services C', false,super_user_id, super_user_id, current_epoch_time, current_epoch_time);

END $$;

-- Create Functions
CREATE OR REPLACE FUNCTION prevent_service_deletion()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE services
    SET is_active = false,
        updated_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)
    WHERE id = OLD.id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create Triggers
CREATE TRIGGER tr_prevent_service_deletion
    BEFORE DELETE ON services
    FOR EACH ROW
EXECUTE FUNCTION prevent_service_deletion();