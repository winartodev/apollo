-- Remove Triggers
DROP TRIGGER IF EXISTS tr_prevent_service_deletion ON services;

-- Remove Functions
DROP FUNCTION IF EXISTS  prevent_service_deletion;

-- Remove indexes
DROP INDEX IF EXISTS idx_services_updated_by;
DROP INDEX IF EXISTS idx_services_created_by;
DROP INDEX IF EXISTS idx_services_active_created;
DROP INDEX IF EXISTS idx_services_updated_at;
DROP INDEX IF EXISTS idx_services_created_at;
DROP INDEX IF EXISTS idx_services_is_active;

-- Remove tables
DROP TABLE IF EXISTS services;

-- Remove seeding data
DO $$
    DECLARE super_user_id INT;
BEGIN
    SELECT id INTO super_user_id
    FROM users
    WHERE email = 'super-user@gmail.com';

    DELETE FROM users WHERE id = super_user_id;
END$$
