-- Remove Index
DROP INDEX IF EXISTS idx_user_applications_application_id;

DROP INDEX IF EXISTS idx_application_services_created_at;
DROP INDEX IF EXISTS idx_application_services_is_active;
DROP INDEX IF EXISTS idx_application_services_scope;
DROP INDEX IF EXISTS idx_application_services_application_id;

DROP INDEX IF EXISTS idx_applications_created_at;
DROP INDEX IF EXISTS idx_applications_is_active;

-- Remove Table
DROP TABLE IF EXISTS user_applications;
DROP TABLE IF EXISTS application_services;
DROP TABLE IF EXISTS applications;
