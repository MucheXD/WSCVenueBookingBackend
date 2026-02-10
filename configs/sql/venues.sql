DROP TABLE IF EXISTS venues;

CREATE TABLE venues (
    venue_id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    create_time TIMESTAMP COMMENT 'Create Time',
    name_text VARCHAR(255) COMMENT 'Venue Name',
    location_building_id INTEGER COMMENT 'Building ID',
    type_id INTEGER COMMENT 'Venue Type ID',
    capacity INTEGER COMMENT 'Venue Capacity',
    description_text TEXT COMMENT 'Venue Description',
    cover_image_file VARCHAR(64) COMMENT 'Cover Image File Hash',
    images_file_list TEXT COMMENT 'List of Image File Hashes',
    is_active BOOLEAN COMMENT 'Is Active Flag',
    FOREIGN KEY (location_building_id) REFERENCES venue_buildings (building_id),
    FOREIGN KEY (cover_image_file) REFERENCES file_objects (file_hash),
    Foreign Key (type_id) REFERENCES venue_types (venue_type_id)
) COMMENT 'Venue List Table';