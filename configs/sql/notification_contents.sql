DROP TABLE IF EXISTS notification_content;

CREATE TABLE notification_content (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each notification content',
    type INT NOT NULL COMMENT 'Type of the notification content (use enum e.g., 1="system", 2="application_result")',
    sender_uid CHAR(12) NOT NULL COMMENT 'User id of the sender, null for system notifications',
    title_text VARCHAR(255) NOT NULL COMMENT 'Title of the notification',
    content_text TEXT COMMENT 'Body content of the notification',
    payload JSON COMMENT 'Additional data for the notification in JSON format',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation timestamp of the notification',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update timestamp of the notification',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp'
) COMMENT 'Table for storing notification content';