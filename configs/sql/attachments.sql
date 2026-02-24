DROP TABLE IF EXISTS attachments;

CREATE TABLE attachments (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each attachment',
    biz_type INT NOT NULL COMMENT 'ID of the associated business type',
    biz_id INT NOT NULL COMMENT 'ID of the associated business entity',
    biz_index INT COMMENT 'Index of the attachment for the same business entity, starting from 0',
    biz_filetype VARCHAR(32) COMMENT 'Type of the file in biz context (e.g., "image")',
    biz_filename VARCHAR(255) COMMENT 'Original filename of the attachment',
    file_token char(64) NOT NULL COMMENT 'File token of the attachment file',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp',
    Foreign Key (file_token) REFERENCES file_objects (file_token)
) COMMENT 'Table for storing attachments related to venue applications';