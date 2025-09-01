-- Add RBAC-related columns to existing users table if they don't exist
-- This migration ensures compatibility with existing user table structure

-- Check and add role-related columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_admin TINYINT(1) DEFAULT 0 AFTER password_hash,
ADD COLUMN IF NOT EXISTS school_id CHAR(36) DEFAULT NULL AFTER is_admin,
ADD COLUMN IF NOT EXISTS majority_id CHAR(36) DEFAULT NULL AFTER school_id,
ADD COLUMN IF NOT EXISTS class_id CHAR(36) DEFAULT NULL AFTER majority_id,
ADD COLUMN IF NOT EXISTS partner_id CHAR(36) DEFAULT NULL AFTER class_id;

-- Add indexes for new columns if they don't exist
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);
CREATE INDEX IF NOT EXISTS idx_users_school_id ON users(school_id);
CREATE INDEX IF NOT EXISTS idx_users_majority_id ON users(majority_id);
CREATE INDEX IF NOT EXISTS idx_users_class_id ON users(class_id);
CREATE INDEX IF NOT EXISTS idx_users_partner_id ON users(partner_id);

-- Add foreign key constraints to users table if they don't exist
-- Note: These will be added after schools, majorities, classes, and partners tables are created

-- Add audit columns to users table if they don't exist
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS created_by CHAR(36) DEFAULT NULL AFTER created_at,
ADD COLUMN IF NOT EXISTS updated_by CHAR(36) DEFAULT NULL AFTER updated_at,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL AFTER updated_by,
ADD COLUMN IF NOT EXISTS deleted_by CHAR(36) DEFAULT NULL AFTER deleted_at;

-- Add indexes for audit columns
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);
CREATE INDEX IF NOT EXISTS idx_users_updated_by ON users(updated_by);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_by ON users(deleted_by);
