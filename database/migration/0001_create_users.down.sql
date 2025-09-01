-- Drop tables in reverse order to avoid foreign key constraint violations

-- First remove foreign key constraints that reference users table
ALTER TABLE schools DROP FOREIGN KEY IF EXISTS fk_schools_created_by;
ALTER TABLE schools DROP FOREIGN KEY IF EXISTS fk_schools_updated_by;
ALTER TABLE schools DROP FOREIGN KEY IF EXISTS fk_schools_deleted_by;

ALTER TABLE majorities DROP FOREIGN KEY IF EXISTS fk_majorities_created_by;
ALTER TABLE majorities DROP FOREIGN KEY IF EXISTS fk_majorities_updated_by;
ALTER TABLE majorities DROP FOREIGN KEY IF EXISTS fk_majorities_deleted_by;

ALTER TABLE classes DROP FOREIGN KEY IF EXISTS fk_classes_created_by;
ALTER TABLE classes DROP FOREIGN KEY IF EXISTS fk_classes_updated_by;
ALTER TABLE classes DROP FOREIGN KEY IF EXISTS fk_classes_deleted_by;

ALTER TABLE partners DROP FOREIGN KEY IF EXISTS fk_partners_created_by;
ALTER TABLE partners DROP FOREIGN KEY IF EXISTS fk_partners_updated_by;
ALTER TABLE partners DROP FOREIGN KEY IF EXISTS fk_partners_deleted_by;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS otps;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS partners;
DROP TABLE IF EXISTS majorities;
DROP TABLE IF EXISTS schools;