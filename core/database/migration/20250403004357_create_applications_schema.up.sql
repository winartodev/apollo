-- Create Table
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT '',
    is_active bool DEFAULT TRUE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS application_services (
    id SERIAL PRIMARY KEY,
    application_id INT REFERENCES applications(id) NOT NULL,
    scope INT DEFAULT 0, -- 0: Public 1:Protected 2:Internal .etc
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    is_active bool DEFAULT TRUE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_applications (
    user_id INT REFERENCES users(id) NOT NULL,
    application_id INT REFERENCES applications(id) NOT NULL,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    PRIMARY KEY(user_id, application_id)
);

-- Create Index
CREATE INDEX idx_applications_is_active ON applications(is_active);
CREATE INDEX idx_applications_created_at ON applications(created_at);

CREATE INDEX idx_application_services_application_id ON application_services(application_id);
CREATE INDEX idx_application_services_scope ON application_services(scope);
CREATE INDEX idx_application_services_is_active ON application_services(is_active);
CREATE INDEX idx_application_services_created_at ON application_services(created_at);

CREATE INDEX idx_user_applications_application_id ON user_applications(application_id);
-- Data Seeding
DO $$
    DECLARE current_epoch_time BIGINT;
BEGIN
    current_epoch_time =  EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);

    INSERT INTO applications (slug, name, description, created_at, updated_at)
    VALUES ('apollo', 'Apollo', 'Provides public user access to the platform', current_epoch_time, current_epoch_time),
           ('apollo-internal', 'Apollo Internal', 'Internal application for employees and administrators',current_epoch_time, current_epoch_time),
           ('apollo-protected', 'Apollo Protected', 'Public application with additional features for registered users',current_epoch_time, current_epoch_time);

    INSERT INTO application_services(application_id, scope, slug, name, is_active, created_at, updated_at)
    VALUES (2, 2, 'test-internal-services-1', 'Test Internal Services 1', true, current_epoch_time, current_epoch_time),
           (2, 2, 'test-internal-services-2', 'Test Internal Services 2', true, current_epoch_time, current_epoch_time),
           (2, 2, 'test-internal-services-3', 'Test Internal Services 3', true, current_epoch_time, current_epoch_time);

    INSERT INTO user_applications(user_id, application_id, created_at, updated_at)
    VALUES (1, 1, current_epoch_time, current_epoch_time),
           (2, 2, current_epoch_time, current_epoch_time),
           (3, 3, current_epoch_time, current_epoch_time);
END$$
