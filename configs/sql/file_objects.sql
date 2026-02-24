DROP TABLE IF EXISTS file_objects;

CREATE TABLE file_objects (
    fid int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '(DBOnly) Primary Key, unique file identifier',
    file_token char(64) NOT NULL UNIQUE COMMENT 'Unique token for file access, calculated from hash and salt',
    file_hash char(64) NOT NULL UNIQUE COMMENT 'Hash of the file',
    file_size bigint COMMENT 'Size of the file in bytes',
    storage_type int NOT NULL COMMENT 'Storage type (e.g., 0 for local, 1 for OSS)',
    link_count int DEFAULT 0 COMMENT '(RepoLayerOnly) Number of references to this file, used for garbage collection'
) COMMENT 'File Objects Table';