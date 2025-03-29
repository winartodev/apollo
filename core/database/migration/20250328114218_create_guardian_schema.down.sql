-- Remove Triggers
DROP TRIGGER IF EXISTS trigger_prevent_guardian_permissions_deletions ON guardian_permissions;
DROP TRIGGER IF EXISTS trigger_prevent_guardian_groups_deletions ON guardian_groups;
DROP TRIGGER IF EXISTS trigger_prevent_guardian_roles_deletions ON guardian_roles;

-- Remove Functions
DROP FUNCTION IF EXISTS trigger_prevent_guardian_permissions_deletions();
DROP FUNCTION IF EXISTS trigger_prevent_guardian_groups_deletions();
DROP FUNCTION IF EXISTS trigger_prevent_guardian_roles_deletions();

-- Remove Indexes
DROP INDEX IF EXISTS idx_guardian_user_applications_is_active;

-- Remove Relationship Tables
DROP TABLE IF EXISTS guardian_user_applications;

-- Remove Tables
DROP TABLE IF EXISTS guardian_permissions;
DROP TABLE IF EXISTS guardian_roles;
DROP TABLE IF EXISTS guardian_groups;
DROP TABLE IF EXISTS applications;
