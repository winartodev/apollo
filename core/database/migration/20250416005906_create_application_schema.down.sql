DROP TRIGGER IF EXISTS  tr_prevent_application_deletion ON applications;

DROP FUNCTION IF EXISTS prevent_application_deletion;

DROP INDEX IF EXISTS idx_services_application_id;

DROP INDEX IF EXISTS idx_application_scope_created_at;
DROP INDEX IF EXISTS idx_application_scope_id;

DROP INDEX IF EXISTS idx_scopes_created_at;
DROP INDEX IF EXISTS idx_scopes_id;

DROP INDEX IF EXISTS idx_applications_updated_by;
DROP INDEX IF EXISTS idx_applications_created_by;
DROP INDEX IF EXISTS idx_applications_active_created;
DROP INDEX IF EXISTS idx_applications_updated_at;
DROP INDEX IF EXISTS idx_applications_created_at;
DROP INDEX IF EXISTS idx_applications_is_active;

ALTER TABLE services DROP COLUMN IF EXISTS application_id;

DROP TABLE IF EXISTS user_applications;
DROP TABLE IF EXISTS application_scope;
DROP TABLE IF EXISTS scopes;
DROP TABLE IF EXISTS applications;
