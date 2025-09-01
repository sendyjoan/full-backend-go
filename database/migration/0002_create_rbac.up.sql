-- Create roles table
CREATE TABLE roles (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_roles_name (name),
  INDEX idx_roles_slug (slug),
  INDEX idx_roles_is_active (is_active),
  INDEX idx_roles_created_by (created_by),
  INDEX idx_roles_deleted_at (deleted_at)
);

-- Create permissions table
CREATE TABLE permissions (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  resource VARCHAR(100) NOT NULL,
  action VARCHAR(50) NOT NULL,
  description TEXT,
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_permissions_name (name),
  INDEX idx_permissions_slug (slug),
  INDEX idx_permissions_resource (resource),
  INDEX idx_permissions_action (action),
  INDEX idx_permissions_resource_action (resource, action),
  INDEX idx_permissions_is_active (is_active),
  INDEX idx_permissions_created_by (created_by),
  INDEX idx_permissions_deleted_at (deleted_at)
);

-- Create menus table
CREATE TABLE menus (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  url VARCHAR(255),
  icon VARCHAR(100),
  parent_id CHAR(36) DEFAULT NULL,
  sort_order INT DEFAULT 0,
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_menus_name (name),
  INDEX idx_menus_slug (slug),
  INDEX idx_menus_parent_id (parent_id),
  INDEX idx_menus_sort_order (sort_order),
  INDEX idx_menus_is_active (is_active),
  INDEX idx_menus_created_by (created_by),
  INDEX idx_menus_deleted_at (deleted_at),
  
  -- Foreign key constraints
  FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE SET NULL
);

-- Create role_permissions junction table (many-to-many)
CREATE TABLE role_permissions (
  id CHAR(36) PRIMARY KEY,
  role_id CHAR(36) NOT NULL,
  permission_id CHAR(36) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_role_permissions_role_id (role_id),
  INDEX idx_role_permissions_permission_id (permission_id),
  INDEX idx_role_permissions_created_by (created_by),
  UNIQUE KEY unique_role_permission (role_id, permission_id),
  
  -- Foreign key constraints
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Create user_roles junction table (many-to-many)
CREATE TABLE user_roles (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  role_id CHAR(36) NOT NULL,
  assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  assigned_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_user_roles_user_id (user_id),
  INDEX idx_user_roles_role_id (role_id),
  INDEX idx_user_roles_assigned_by (assigned_by),
  UNIQUE KEY unique_user_role (user_id, role_id),
  
  -- Foreign key constraints
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Create role_menus junction table (many-to-many)
CREATE TABLE role_menus (
  id CHAR(36) PRIMARY KEY,
  role_id CHAR(36) NOT NULL,
  menu_id CHAR(36) NOT NULL,
  can_view TINYINT(1) DEFAULT 1,
  can_create TINYINT(1) DEFAULT 0,
  can_edit TINYINT(1) DEFAULT 0,
  can_delete TINYINT(1) DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_role_menus_role_id (role_id),
  INDEX idx_role_menus_menu_id (menu_id),
  INDEX idx_role_menus_created_by (created_by),
  UNIQUE KEY unique_role_menu (role_id, menu_id),
  
  -- Foreign key constraints
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
);

-- Add foreign key constraints after all tables are created
ALTER TABLE roles ADD CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE roles ADD CONSTRAINT fk_roles_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE roles ADD CONSTRAINT fk_roles_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

ALTER TABLE permissions ADD CONSTRAINT fk_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

ALTER TABLE menus ADD CONSTRAINT fk_menus_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE menus ADD CONSTRAINT fk_menus_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE menus ADD CONSTRAINT fk_menus_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id);
ALTER TABLE role_menus ADD CONSTRAINT fk_role_menus_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE role_menus ADD CONSTRAINT fk_role_menus_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);

-- Insert default roles
INSERT INTO roles (id, name, slug, description, is_active, created_at) VALUES
(UUID(), 'Super Admin', 'super-admin', 'Super administrator with full access to all system features', 1, NOW()),
(UUID(), 'Admin', 'admin', 'Administrator with access to most system features', 1, NOW()),
(UUID(), 'School Admin', 'school-admin', 'School administrator with access to school-specific features', 1, NOW()),
(UUID(), 'Teacher', 'teacher', 'Teacher with access to academic features', 1, NOW()),
(UUID(), 'Student', 'student', 'Student with limited access to academic features', 1, NOW()),
(UUID(), 'Partner', 'partner', 'Partner with access to partnership-related features', 1, NOW());

-- Insert default permissions
INSERT INTO permissions (id, name, slug, resource, action, description, is_active, created_at) VALUES
-- User management permissions
(UUID(), 'View Users', 'view-users', 'users', 'view', 'Permission to view user list and details', 1, NOW()),
(UUID(), 'Create Users', 'create-users', 'users', 'create', 'Permission to create new users', 1, NOW()),
(UUID(), 'Edit Users', 'edit-users', 'users', 'edit', 'Permission to edit existing users', 1, NOW()),
(UUID(), 'Delete Users', 'delete-users', 'users', 'delete', 'Permission to delete users', 1, NOW()),

-- Role management permissions
(UUID(), 'View Roles', 'view-roles', 'roles', 'view', 'Permission to view role list and details', 1, NOW()),
(UUID(), 'Create Roles', 'create-roles', 'roles', 'create', 'Permission to create new roles', 1, NOW()),
(UUID(), 'Edit Roles', 'edit-roles', 'roles', 'edit', 'Permission to edit existing roles', 1, NOW()),
(UUID(), 'Delete Roles', 'delete-roles', 'roles', 'delete', 'Permission to delete roles', 1, NOW()),

-- Permission management permissions
(UUID(), 'View Permissions', 'view-permissions', 'permissions', 'view', 'Permission to view permission list and details', 1, NOW()),
(UUID(), 'Create Permissions', 'create-permissions', 'permissions', 'create', 'Permission to create new permissions', 1, NOW()),
(UUID(), 'Edit Permissions', 'edit-permissions', 'permissions', 'edit', 'Permission to edit existing permissions', 1, NOW()),
(UUID(), 'Delete Permissions', 'delete-permissions', 'permissions', 'delete', 'Permission to delete permissions', 1, NOW()),

-- Menu management permissions
(UUID(), 'View Menus', 'view-menus', 'menus', 'view', 'Permission to view menu list and details', 1, NOW()),
(UUID(), 'Create Menus', 'create-menus', 'menus', 'create', 'Permission to create new menus', 1, NOW()),
(UUID(), 'Edit Menus', 'edit-menus', 'menus', 'edit', 'Permission to edit existing menus', 1, NOW()),
(UUID(), 'Delete Menus', 'delete-menus', 'menus', 'delete', 'Permission to delete menus', 1, NOW()),

-- School management permissions
(UUID(), 'View Schools', 'view-schools', 'schools', 'view', 'Permission to view school list and details', 1, NOW()),
(UUID(), 'Create Schools', 'create-schools', 'schools', 'create', 'Permission to create new schools', 1, NOW()),
(UUID(), 'Edit Schools', 'edit-schools', 'schools', 'edit', 'Permission to edit existing schools', 1, NOW()),
(UUID(), 'Delete Schools', 'delete-schools', 'schools', 'delete', 'Permission to delete schools', 1, NOW());

-- Insert default menus
INSERT INTO menus (id, name, slug, url, icon, parent_id, sort_order, is_active, created_at) VALUES
-- Main menus
(UUID(), 'Dashboard', 'dashboard', '/dashboard', 'dashboard', NULL, 1, 1, NOW()),
(UUID(), 'User Management', 'user-management', '#', 'users', NULL, 2, 1, NOW()),
(UUID(), 'Role & Permissions', 'role-permissions', '#', 'shield', NULL, 3, 1, NOW()),
(UUID(), 'School Management', 'school-management', '#', 'school', NULL, 4, 1, NOW()),
(UUID(), 'Academic', 'academic', '#', 'book', NULL, 5, 1, NOW()),
(UUID(), 'Partnership', 'partnership', '#', 'handshake', NULL, 6, 1, NOW()),

-- User Management submenus
(UUID(), 'Users', 'users', '/users', 'user', (SELECT id FROM menus WHERE slug = 'user-management'), 1, 1, NOW()),
(UUID(), 'User Roles', 'user-roles', '/user-roles', 'user-check', (SELECT id FROM menus WHERE slug = 'user-management'), 2, 1, NOW()),

-- Role & Permissions submenus
(UUID(), 'Roles', 'roles', '/roles', 'shield-check', (SELECT id FROM menus WHERE slug = 'role-permissions'), 1, 1, NOW()),
(UUID(), 'Permissions', 'permissions', '/permissions', 'key', (SELECT id FROM menus WHERE slug = 'role-permissions'), 2, 1, NOW()),
(UUID(), 'Menu Management', 'menu-management', '/menus', 'menu', (SELECT id FROM menus WHERE slug = 'role-permissions'), 3, 1, NOW()),

-- School Management submenus
(UUID(), 'Schools', 'schools', '/schools', 'building', (SELECT id FROM menus WHERE slug = 'school-management'), 1, 1, NOW()),
(UUID(), 'Majorities', 'majorities', '/majorities', 'graduation-cap', (SELECT id FROM menus WHERE slug = 'school-management'), 2, 1, NOW()),
(UUID(), 'Classes', 'classes', '/classes', 'users', (SELECT id FROM menus WHERE slug = 'school-management'), 3, 1, NOW()),

-- Academic submenus
(UUID(), 'Curriculum', 'curriculum', '/curriculum', 'book-open', (SELECT id FROM menus WHERE slug = 'academic'), 1, 1, NOW()),
(UUID(), 'Courses', 'courses', '/courses', 'file-text', (SELECT id FROM menus WHERE slug = 'academic'), 2, 1, NOW()),

-- Partnership submenus
(UUID(), 'Partners', 'partners', '/partners', 'briefcase', (SELECT id FROM menus WHERE slug = 'partnership'), 1, 1, NOW()),
(UUID(), 'Internships', 'internships', '/internships', 'clipboard', (SELECT id FROM menus WHERE slug = 'partnership'), 2, 1, NOW());
