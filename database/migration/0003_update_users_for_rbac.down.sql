-- Remove RBAC-related columns from users table

-- Drop foreign key constraints first if they exist
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_school_id;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_majority_id;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_class_id;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_partner_id;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_created_by;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_updated_by;
ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_users_deleted_by;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_is_admin ON users;
DROP INDEX IF EXISTS idx_users_school_id ON users;
DROP INDEX IF EXISTS idx_users_majority_id ON users;
DROP INDEX IF EXISTS idx_users_class_id ON users;
DROP INDEX IF EXISTS idx_users_partner_id ON users;
DROP INDEX IF EXISTS idx_users_created_by ON users;
DROP INDEX IF EXISTS idx_users_updated_by ON users;
DROP INDEX IF EXISTS idx_users_deleted_at ON users;
DROP INDEX IF EXISTS idx_users_deleted_by ON users;

-- Drop columns
ALTER TABLE users 
DROP COLUMN IF EXISTS is_admin,
DROP COLUMN IF EXISTS school_id,
DROP COLUMN IF EXISTS majority_id,
DROP COLUMN IF EXISTS class_id,
DROP COLUMN IF EXISTS partner_id,
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;
