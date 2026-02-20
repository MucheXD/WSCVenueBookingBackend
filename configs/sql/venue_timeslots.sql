DROP TABLE IF EXISTS venue_timeslots;

CREATE TABLE venue_timeslots (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each timetable info',
    venue_id int not null COMMENT 'Target venue id of this timetable info',
    start_time TIMESTAMP not null COMMENT 'Occupation start time',
    end_time TIMESTAMP COMMENT 'Occupation end time, null for not sure',
    application_id int COMMENT 'Associated application id, null if not applicable',
    status char(8) not null COMMENT 'Timeslot status',
    Foreign Key (venue_id) REFERENCES venues (venue_id),
    Foreign Key (application_id) REFERENCES applications (id)
)