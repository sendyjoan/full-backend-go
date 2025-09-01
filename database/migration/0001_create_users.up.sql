-- Create schools table first (no dependencies)
CREATE TABLE schools (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  address VARCHAR(255),
  domain VARCHAR(255) UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_schools_domain (domain),
  INDEX idx_schools_created_by (created_by),
  INDEX idx_schools_deleted_at (deleted_at)
);

-- Create majorities table (depends on schools)
CREATE TABLE majorities (
  id CHAR(36) PRIMARY KEY,
  school_id CHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_majorities_school_id (school_id),
  INDEX idx_majorities_created_by (created_by),
  INDEX idx_majorities_deleted_at (deleted_at),
  
  -- Foreign key constraints
  FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE
);

-- Create partners table (depends on schools)
CREATE TABLE partners (
  id CHAR(36) PRIMARY KEY,
  school_id CHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  website VARCHAR(255),
  description TEXT,
  address VARCHAR(255),
  contact_name VARCHAR(255),
  contact_person VARCHAR(255),
  contact_email VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_partners_school_id (school_id),
  INDEX idx_partners_created_by (created_by),
  INDEX idx_partners_deleted_at (deleted_at),
  INDEX idx_partners_name (name),
  
  -- Foreign key constraints
  FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE
);

-- Create classes table (depends on schools and majorities)
CREATE TABLE classes (
  id CHAR(36) PRIMARY KEY,
  school_id CHAR(36) NOT NULL,
  majority_id CHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL,
  deleted_by CHAR(36),
  
  -- Indexes for performance
  INDEX idx_classes_school_id (school_id),
  INDEX idx_classes_majority_id (majority_id),
  INDEX idx_classes_created_by (created_by),
  INDEX idx_classes_deleted_at (deleted_at),
  INDEX idx_classes_school_majority (school_id, majority_id), -- Composite index for common queries
  
  -- Foreign key constraints
  FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE,
  FOREIGN KEY (majority_id) REFERENCES majorities(id) ON DELETE CASCADE
);

-- Create users table (depends on schools, majorities, classes, partners)
CREATE TABLE users (
  id CHAR(36) PRIMARY KEY,
  username VARCHAR(60) UNIQUE NOT NULL,
  email VARCHAR(120) UNIQUE NOT NULL,
  fullname VARCHAR(120) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  is_admin TINYINT(1) DEFAULT 0,
  school_id CHAR(36) DEFAULT NULL,
  majority_id CHAR(36) DEFAULT NULL,
  class_id CHAR(36) DEFAULT NULL,
  partner_id CHAR(36) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  -- Indexes for performance
  INDEX idx_users_username (username),
  INDEX idx_users_email (email),
  INDEX idx_users_school_id (school_id),
  INDEX idx_users_majority_id (majority_id),
  INDEX idx_users_class_id (class_id),
  INDEX idx_users_partner_id (partner_id),
  INDEX idx_users_is_admin (is_admin),
  INDEX idx_users_fullname (fullname),
  
  -- Foreign key constraints
  FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE SET NULL,
  FOREIGN KEY (majority_id) REFERENCES majorities(id) ON DELETE SET NULL,
  FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE SET NULL,
  FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE SET NULL
);

-- Create otps table (depends on users)
CREATE TABLE otps (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  code VARCHAR(6) NOT NULL,
  purpose VARCHAR(32) NOT NULL,           -- e.g. "forgot_password"
  expires_at TIMESTAMP NOT NULL,
  used TINYINT(1) DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Indexes for performance
  INDEX idx_otps_user_id (user_id),
  INDEX idx_otps_code (code),
  INDEX idx_otps_purpose (purpose),
  INDEX idx_otps_expires_at (expires_at),
  INDEX idx_otps_used (used),
  INDEX idx_otps_user_purpose (user_id, purpose), -- Composite index for finding user's OTP by purpose
  
  -- Foreign key constraints
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create refresh_tokens table (depends on users)
CREATE TABLE refresh_tokens (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  token_hash VARCHAR(255) NOT NULL,
  user_agent VARCHAR(255),
  ip VARCHAR(64),
  revoked TINYINT(1) DEFAULT 0,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Indexes for performance
  INDEX idx_refresh_tokens_user_id (user_id),
  INDEX idx_refresh_tokens_token_hash (token_hash),
  INDEX idx_refresh_tokens_expires_at (expires_at),
  INDEX idx_refresh_tokens_revoked (revoked),
  INDEX idx_refresh_tokens_user_active (user_id, revoked, expires_at), -- Composite index for finding active tokens
  
  -- Foreign key constraints
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add foreign key constraints to schools table for created_by, updated_by, deleted_by after users table is created
ALTER TABLE schools ADD CONSTRAINT fk_schools_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE schools ADD CONSTRAINT fk_schools_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE schools ADD CONSTRAINT fk_schools_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- Add foreign key constraints to majorities table for created_by, updated_by, deleted_by
ALTER TABLE majorities ADD CONSTRAINT fk_majorities_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE majorities ADD CONSTRAINT fk_majorities_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE majorities ADD CONSTRAINT fk_majorities_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- Add foreign key constraints to classes table for created_by, updated_by, deleted_by
ALTER TABLE classes ADD CONSTRAINT fk_classes_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE classes ADD CONSTRAINT fk_classes_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE classes ADD CONSTRAINT fk_classes_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- Add foreign key constraints to partners table for created_by, updated_by, deleted_by
ALTER TABLE partners ADD CONSTRAINT fk_partners_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE partners ADD CONSTRAINT fk_partners_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE partners ADD CONSTRAINT fk_partners_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);