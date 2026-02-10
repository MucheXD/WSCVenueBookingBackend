DROP TABLE IF EXISTS venue_types;

CREATE TABLE venue_types (
    venue_type_id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    name_text VARCHAR(255) COMMENT 'Venue Type Name',
    description_text TEXT COMMENT 'Venue Type Description'
) COMMENT 'Venue Type Table';