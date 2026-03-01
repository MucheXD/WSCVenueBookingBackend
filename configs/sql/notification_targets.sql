DROP TABLE IF EXISTS notification_targets;

CREATE TABLE notification_targets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each notification target',
    type INT NOT NULL COMMENT 'Type of the notification content (use enum e.g., 1="system", 2="application_result")',
    notification_id INT NOT NULL COMMENT 'ID of the associated notification content',
    receiver_uid CHAR(12) NOT NULL COMMENT 'User id of the receiver',
    is_read BOOLEAN DEFAULT FALSE COMMENT 'True if the notification was read',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp, users can delete to hide the notification for themselves',
    UNIQUE KEY uk_notification_receiver (notification_id, receiver_uid),
    Foreign Key (notification_id) REFERENCES notification_contents (id)
) COMMENT 'Table for storing notification targets and their read status';