-- Drop tables in reverse order to avoid foreign key constraint violations

-- Drop foreign key constraints first
ALTER TABLE role_menus DROP FOREIGN KEY IF EXISTS fk_role_menus_created_by;
ALTER TABLE role_menus DROP FOREIGN KEY IF EXISTS fk_role_menus_updated_by;
ALTER TABLE user_roles DROP FOREIGN KEY IF EXISTS fk_user_roles_assigned_by;
ALTER TABLE role_permissions DROP FOREIGN KEY IF EXISTS fk_role_permissions_created_by;
ALTER TABLE menus DROP FOREIGN KEY IF EXISTS fk_menus_created_by;
ALTER TABLE menus DROP FOREIGN KEY IF EXISTS fk_menus_updated_by;
ALTER TABLE menus DROP FOREIGN KEY IF EXISTS fk_menus_deleted_by;
ALTER TABLE permissions DROP FOREIGN KEY IF EXISTS fk_permissions_created_by;
ALTER TABLE permissions DROP FOREIGN KEY IF EXISTS fk_permissions_updated_by;
ALTER TABLE permissions DROP FOREIGN KEY IF EXISTS fk_permissions_deleted_by;
ALTER TABLE roles DROP FOREIGN KEY IF EXISTS fk_roles_created_by;
ALTER TABLE roles DROP FOREIGN KEY IF EXISTS fk_roles_updated_by;
ALTER TABLE roles DROP FOREIGN KEY IF EXISTS fk_roles_deleted_by;

-- Drop junction tables first
DROP TABLE IF EXISTS role_menus;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;

-- Drop main tables
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
