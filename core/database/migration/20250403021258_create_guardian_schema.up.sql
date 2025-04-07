-- Create Tables
CREATE TABLE IF NOT EXISTS guardian_roles (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT '',
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guardian_permissions (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT '',
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guardian_user_roles (
    user_id INT REFERENCES users(id) NOT NULL,
    application_id INT REFERENCES applications(id) NOT NULL,
    role_id INT REFERENCES guardian_roles(id) NOT NULL,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    PRIMARY KEY (user_id, application_id)
);

CREATE TABLE IF NOT EXISTS guardian_role_permissions (
    role_id INT REFERENCES guardian_roles(id) NOT NULL,
    application_service_id INT REFERENCES application_services(id) NOT NULL,
    permission_id INT REFERENCES guardian_permissions(id) NOT NULL,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    PRIMARY KEY (role_id, application_service_id, permission_id)
);

-- Indexing
CREATE INDEX IF NOT EXISTS idx_guardian_user_roles_role_id ON guardian_user_roles(role_id);

-- Data Seeding
DO $$
    DECLARE current_epoch_time BIGINT;
BEGIN
    current_epoch_time =  EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);

    INSERT INTO guardian_roles(slug, name, description, created_at, updated_at)
    VALUES ('super-admin', 'Super Administrator', 'Has full access to all system features and settings', current_epoch_time, current_epoch_time),
           ('admin', 'Administrator', 'Can manage most system settings and user permissions', current_epoch_time, current_epoch_time),
           ('moderator', 'Moderator', 'Can manage content and user interactions', current_epoch_time, current_epoch_time),
           ('user', 'Regular User', 'Standard authenticated user with basic privileges', current_epoch_time, current_epoch_time),
           ('guest', 'Guest User', 'Limited access for unauthenticated users', current_epoch_time, current_epoch_time),
           ('developer', 'Developer', 'Technical staff with API and system integration access', current_epoch_time, current_epoch_time),
           ('support', 'Support Staff', 'Can access customer support features and user accounts', current_epoch_time, current_epoch_time),
           ('content-manager', 'Content Manager', 'Can create, edit, and publish content', current_epoch_time, current_epoch_time),
           ('billing', 'Billing Specialist', 'Manages subscriptions, payments, and invoices', current_epoch_time, current_epoch_time),
           ('api-user', 'API User', 'Service account for system integrations', current_epoch_time, current_epoch_time);

    INSERT INTO guardian_permissions(slug, name, description, created_at, updated_at)
    VALUES ('get', 'GET', 'Retrieve Data', current_epoch_time, current_epoch_time),
           ('post', 'POST', 'Write Data', current_epoch_time, current_epoch_time),
           ('put', 'PUT', 'Update Data', current_epoch_time, current_epoch_time),
           ('delete', 'DELETE', 'Delete Data', current_epoch_time, current_epoch_time);

    INSERT INTO guardian_user_roles (user_id, application_id, role_id, created_at, updated_at)
        VALUES (2,2, 2, current_epoch_time, current_epoch_time);

    INSERT INTO guardian_role_permissions (role_id, application_service_id, permission_id, created_at, updated_at)
    VALUES (2,1,1, current_epoch_time, current_epoch_time),
           (2,1,2, current_epoch_time, current_epoch_time),
           (2,1,3, current_epoch_time, current_epoch_time),
           (2,1,4, current_epoch_time, current_epoch_time),
           (2,2,1, current_epoch_time, current_epoch_time),
           (2,2,2, current_epoch_time, current_epoch_time);

END $$