DROP TABLE IF EXISTS users;

CREATE TABLE users (
    uid int NOT NULL PRIMARY KEY COMMENT 'Unique User ID',
    password_hash VARCHAR(64) COMMENT 'Password Hash',
    registered_at TIMESTAMP COMMENT 'Account Creation Time',
    username VARCHAR(255) COMMENT 'Username',
    school_id VARCHAR(255) COMMENT 'School ID, student/staff number',
    real_name VARCHAR(255) COMMENT 'Real Name',
    perm_type char(8) COMMENT 'Permission Level',
    perm_vagid int COMMENT 'Venue Access Group ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record Creation Time',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record Update Time',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT 'Record Deletion Time'
) COMMENT 'User List Table';