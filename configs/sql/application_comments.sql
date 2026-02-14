DROP TABLE IF EXISTS application_comments;

CREATE TABLE application_comments (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each comment',
    application_id INT NOT NULL COMMENT 'ID of the application this comment belongs to',
    commenter_uid INT NOT NULL COMMENT 'User ID of the commenter',
    comment_text TEXT COMMENT 'Text content of the comment',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the comment was created',
    deleted_at TIMESTAMP NULL COMMENT 'Timestamp when the comment was deleted, NULL if not deleted',
    has_attachments BOOLEAN DEFAULT FALSE COMMENT 'Indicates if there are attachments in the comment',
    Foreign Key (application_id) REFERENCES applications (id),
    Foreign Key (commenter_uid) REFERENCES users (uid)
) COMMENT 'Comments on Venue Applications Table';