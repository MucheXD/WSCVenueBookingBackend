DROP TABLE IF EXISTS applications;

CREATE TABLE applications (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each application, NOT for use',
    venue_id int not null COMMENT 'Target venue id of this application',
    applicant_uid CHAR(12) not null COMMENT 'User id of the applicant',
    application_status varchar(16) not null COMMENT 'Status of the application',
    description_text TEXT COMMENT 'Description for application',
    estimated_participants INT COMMENT 'Estimated number of participants',
    has_attachments BOOLEAN DEFAULT FALSE COMMENT '(RepoLayerOnly) Indicates if there are attachments',
    activity_name varchar(255) COMMENT '(ApplicationForm) Name of the activity',
    activity_organizer varchar(255) COMMENT '(ApplicationForm) Organizer of the activity',
    activity_coordinator JSON COMMENT '(ApplicationForm) Coordinator details in JSON format',
    deleted_at TIMESTAMP NULL COMMENT '(RepoLayerOnly) Soft delete timestamp',
    Foreign Key (venue_id) REFERENCES venues (venue_id),
    Foreign Key (applicant_uid) REFERENCES users (uid)
) COMMENT 'Venue Applications Table';