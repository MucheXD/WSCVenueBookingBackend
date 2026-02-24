DROP TABLE IF EXISTS venues;

CREATE TABLE venues (
    venue_id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    create_time TIMESTAMP COMMENT 'Create Time',
    name_text VARCHAR(255) COMMENT 'Venue Name',
    location_building_id INTEGER COMMENT 'Building ID',
    type_id INTEGER COMMENT 'Venue Type ID',
    capacity INTEGER COMMENT 'Venue Capacity',
    description_text TEXT COMMENT 'Venue Description',
    cover_image_token char(64) COMMENT 'Cover Image File Token',
    is_active BOOLEAN COMMENT 'Is Active Flag',
    delete_at TIMESTAMP COMMENT 'Delete Time',
    update_at TIMESTAMP COMMENT 'Update Time',
    FOREIGN KEY (location_building_id) REFERENCES venue_buildings (building_id),
    FOREIGN KEY (cover_image_token) REFERENCES file_objects (file_token),
    Foreign Key (type_id) REFERENCES venue_types (venue_type_id)
) COMMENT 'Venue List Table';