DROP TABLE IF EXISTS venue_buildings;

CREATE TABLE venue_buildings (
    building_id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    building_name varchar(255) NOT NULL COMMENT 'Name of the building',
    location_campus_id int NOT NULL COMMENT 'Foreign Key referencing campus_id in campus_list',
    FOREIGN KEY (location_campus_id) REFERENCES venue_campuses (campus_id)
) COMMENT 'Building List Table';