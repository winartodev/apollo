CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOL DEFAULT TRUE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guardian_groups (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) UNIQUE NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_protected bool DEFAULT TRUE,
    is_deleted bool DEFAULT FALSE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    deleted_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guardian_roles (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) UNIQUE NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    is_protected bool DEFAULT TRUE,
    is_deleted bool DEFAULT FALSE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    deleted_at BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guardian_permissions (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(55) UNIQUE NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_protected bool DEFAULT TRUE,
    is_deleted bool DEFAULT FALSE,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    deleted_at BIGINT DEFAULT 0
);

-- Relationship Tables
CREATE TABLE IF NOT EXISTS guardian_user_applications (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_id INT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    is_active BOOL DEFAULT TRUE,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id),
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    PRIMARY KEY(user_id, application_id)
);

-- Indexing
CREATE INDEX idx_guardian_user_applications_is_active ON guardian_user_applications(is_active) WHERE is_active = TRUE;

-- Data Seeding
INSERT INTO guardian_roles (slug, name, created_at, updated_at)
VALUES ('super-user', 'Super User', EXTRACT(EPOCH FROM CURRENT_TIMESTAMP), EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)),
       ('user', 'User', EXTRACT(EPOCH FROM CURRENT_TIMESTAMP), EXTRACT(EPOCH FROM CURRENT_TIMESTAMP));


-- Functions
CREATE OR REPLACE FUNCTION trigger_prevent_guardian_roles_deletions()
RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.is_protected THEN
        RAISE EXCEPTION 'Cannot delete protected row';
    END IF;

    IF TG_OP = 'DELETE' THEN

        UPDATE guardian_roles SET
            is_deleted = TRUE,
            deleted_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)
        WHERE id = OLD.id;

        RETURN NULL;
    END IF;

    RETURN OLD;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trigger_prevent_guardian_groups_deletions()
RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.is_protected THEN
        RAISE EXCEPTION 'Cannot delete protected row';
    END IF;

    IF TG_OP = 'DELETE' THEN

        UPDATE guardian_groups SET
            is_deleted = TRUE,
            deleted_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)
        WHERE id = OLD.id;

        RETURN NULL;
    END IF;

    RETURN OLD;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trigger_prevent_guardian_permissions_deletions()
RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.is_protected THEN
        RAISE EXCEPTION 'Cannot delete protected row';
    END IF;

    IF TG_OP = 'DELETE' THEN

        UPDATE guardian_permissions SET
            is_deleted = TRUE,
            deleted_at = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)
        WHERE id = OLD.id;

        RETURN NULL;
    END IF;

    RETURN OLD;
END;
$$
LANGUAGE plpgsql;

-- Triggers
CREATE TRIGGER trigger_prevent_guardian_roles_deletions
    BEFORE DELETE ON guardian_roles
    FOR EACH ROW
EXECUTE FUNCTION trigger_prevent_guardian_roles_deletions();

CREATE TRIGGER trigger_prevent_guardian_groups_deletions
    BEFORE DELETE ON guardian_groups
    FOR EACH ROW
EXECUTE FUNCTION trigger_prevent_guardian_groups_deletions();

CREATE TRIGGER trigger_prevent_guardian_permissions_deletions
    BEFORE DELETE ON guardian_permissions
    FOR EACH ROW
EXECUTE FUNCTION trigger_prevent_guardian_permissions_deletions();
