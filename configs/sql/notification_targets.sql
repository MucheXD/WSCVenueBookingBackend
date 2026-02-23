DROP TABLE IF EXISTS notification_targets;

CREATE TABLE notification_targets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each notification target',
    notification_id INT NOT NULL COMMENT 'ID of the associated notification content',
    receiver_uid CHAR(12) NOT NULL COMMENT 'User id of the receiver',
    read_at TIMESTAMP NULL COMMENT 'Timestamp when the notification was read, null if unread',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp, users can delete to hide the notification for themselves',
    Foreign Key (notification_id) REFERENCES notification_content (id)
) COMMENT 'Table for storing notification targets and their read status';