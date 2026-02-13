DROP TABLE IF EXISTS attachments;

CREATE TABLE attachments (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each attachment',
    biz_type INT NOT NULL COMMENT 'ID of the associated business type',
    biz_id INT NOT NULL COMMENT 'ID of the associated business entity',
    fid int NOT NULL COMMENT 'File ID of the attachment file',
    Foreign Key (fid) REFERENCES file_objects (fid)
) COMMENT 'Table for storing attachments related to venue applications';