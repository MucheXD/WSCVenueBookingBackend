DROP TABLE IF EXISTS notifications;

CREATE TABLE notifications (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each notification content',
    sender_uid CHAR(12) NOT NULL COMMENT 'User id of the sender, null for system notifications',
    recevier_uid CHAR(12) NOT NULL COMMENT 'User id of the recevier',
    title VARCHAR(255) NOT NULL COMMENT 'Title of the notification',
    content TEXT COMMENT 'Body content of the notification',
    status INT NOT NULL DEFAULT '0' COMMENT 'Form of the notification',
    release_time VARCHAR(255) NOT NULL COMMENT 'Release time of the notification',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation timestamp of the notification',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update timestamp of the notification',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp'
) COMMENT 'Table for storing notification content';