DROP TABLE IF EXISTS venue_campuses;

CREATE TABLE venue_campuses (
    campus_id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    campus_name varchar(255) NOT NULL COMMENT 'Name of the campus'
) COMMENT 'Campus List Table';