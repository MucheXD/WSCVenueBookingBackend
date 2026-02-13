DROP TABLE IF EXISTS file_objects;

CREATE TABLE file_objects (
    fid int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key, unique file identifier',
    file_hash char(64) NOT NULL UNIQUE COMMENT 'Hash of the file',
    file_name varchar(255) NOT NULL COMMENT 'File name in the storage',
    file_title varchar(255) COMMENT 'Name of the file (Optional)',
    file_type varchar(16) COMMENT 'Type of the file (e.g., media, document, etc.)',
    file_size bigint COMMENT 'Size of the file in bytes'
) COMMENT 'File Objects Table';