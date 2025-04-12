-- Create Tables
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOL NOT NULL DEFAULT TRUE,
    created_by INT NOT NULL REFERENCES users(id),
    updated_by INT NOT NULL REFERENCES users(id),
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS scopes (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS application_scope (
    id SERIAL UNIQUE,
    application_id INT NOT NULL REFERENCES applications(id) ON DELETE RESTRICT,
    scope_id INT NOT NULL REFERENCES scopes(id) ON DELETE RESTRICT,
    is_active BOOL NOT NULL DEFAULT TRUE,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    PRIMARY KEY (application_id, scope_id)
);

CREATE TABLE IF NOT EXISTS user_applications (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_scope_id INT NOT NULL REFERENCES application_scope(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, application_scope_id)
);

-- Create ALTER TABLE
ALTER TABLE services ADD COLUMN application_id INT REFERENCES applications(id) ON DELETE RESTRICT;

-- Created Indexes
CREATE INDEX IF NOT EXISTS idx_applications_is_active ON applications(is_active);
CREATE INDEX IF NOT EXISTS idx_applications_created_at ON applications(created_at);
CREATE INDEX IF NOT EXISTS idx_applications_updated_at ON applications(updated_at);
CREATE INDEX IF NOT EXISTS idx_applications_active_created ON applications(is_active, created_at);
CREATE INDEX IF NOT EXISTS idx_applications_created_by ON applications(created_by);
CREATE INDEX IF NOT EXISTS idx_applications_updated_by ON applications(updated_by);

CREATE INDEX IF NOT EXISTS idx_scopes_id ON scopes(id);
CREATE INDEX IF NOT EXISTS idx_scopes_created_at ON scopes(created_at);

CREATE INDEX IF NOT EXISTS idx_application_scope_id ON application_scope(id);
CREATE INDEX IF NOT EXISTS idx_application_scope_created_at ON application_scope(created_at);

CREATE INDEX IF NOT EXISTS idx_services_application_id ON services(application_id);

-- Seeding Data
DO $$
    DECLARE current_epoch_time BIGINT;
    DECLARE super_user_id INT;
    DECLARE app_id INT;
BEGIN
    current_epoch_time =  EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);

    SELECT id INTO super_user_id FROM users WHERE email = 'super-user@gmail.com';

    INSERT INTO scopes (slug, name, created_at, updated_at)
    VALUES ('public', 'Public', current_epoch_time, current_epoch_time),
    ('internal', 'Internal', current_epoch_time, current_epoch_time),
    ('protected', 'Protected', current_epoch_time, current_epoch_time);

    INSERT INTO applications (slug, name, is_active, created_by, updated_by, created_at, updated_at)
    VALUES ('apollo', 'Apollo', true, super_user_id, super_user_id, current_epoch_time, current_epoch_time)
    RETURNING id INTO app_id;

    INSERT INTO application_scope (application_id, scope_id, is_active, created_by, updated_by, created_at, updated_at)
    VALUES (app_id, (SELECT id AS scope_id FROM scopes WHERE slug = 'public'), true, super_user_id, super_user_id, current_epoch_time, current_epoch_time),
           (app_id, (SELECT id AS scope_id FROM scopes WHERE slug = 'internal'), true, super_user_id, super_user_id, current_epoch_time, current_epoch_time),
           (app_id, (SELECT id AS scope_id FROM scopes WHERE slug = 'protected'), true, super_user_id, super_user_id, current_epoch_time, current_epoch_time);

    INSERT INTO user_applications (user_id, application_scope_id)
    VALUES (super_user_id, (SELECT app_scope.id AS application_scope_id FROM application_scope AS app_scope JOIN scopes AS s ON app_scope.scope_id = s.id WHERE s.slug = 'public')),
           (super_user_id, (SELECT app_scope.id AS application_scope_id FROM application_scope AS app_scope JOIN scopes AS s ON app_scope.scope_id = s.id WHERE s.slug = 'internal')),
           (super_user_id, (SELECT app_scope.id AS application_scope_id FROM application_scope AS app_scope JOIN scopes AS s ON app_scope.scope_id = s.id WHERE s.slug = 'protected'));

    UPDATE services SET application_id = app_id WHERE id = 1;
    UPDATE services SET application_id = app_id WHERE id = 2;
    UPDATE services SET application_id = app_id WHERE id = 3;
END $$;

-- Create Functions
CREATE OR REPLACE FUNCTION prevent_application_deletion()
    RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Deletion from applications table is not allowed. To deactivate an application, set is_active to false instead of deleting the record.';
END;
$$ LANGUAGE plpgsql;

-- Create Triggers
CREATE TRIGGER tr_prevent_application_deletion
    BEFORE DELETE ON applications
    FOR EACH ROW
EXECUTE FUNCTION prevent_application_deletion();
