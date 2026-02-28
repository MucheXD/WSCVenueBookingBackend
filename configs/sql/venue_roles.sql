DROP TABLE IF EXISTS venue_roles;

CREATE TABLE venue_roles (
    vagid int NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary Key',
    role_name varchar(255) NOT NULL COMMENT 'Name of the role',
    role_description text COMMENT 'Description of the role'
) COMMENT 'Venue Roles Table';