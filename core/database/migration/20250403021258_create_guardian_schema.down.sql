-- Remove Index
DROP INDEX IF EXISTS idx_guardian_user_roles_role_id;

-- Remove Tables
DROP TABLE IF EXISTS guardian_role_permissions;
DROP TABLE IF EXISTS guardian_user_roles;
DROP TABLE IF EXISTS guardian_permissions;
DROP TABLE IF EXISTS guardian_roles;
