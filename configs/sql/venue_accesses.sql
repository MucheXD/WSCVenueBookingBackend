DROP TABLE IF EXISTS venue_accesses;

CREATE TABLE venue_accesses (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each access info, NOT for use',
    vagid int not null COMMENT 'Unique identifier for each access group',
    venue_id int not null COMMENT 'Target venue id of this access info',
    allow_reservation tinyint(1) not null COMMENT 'Whether this access group allows reservation',
    allow_approval tinyint(1) not null COMMENT 'Whether this access group allows approval',
    allow_edit tinyint(1) not null COMMENT 'Whether this access group allows edit',
    allow_manage tinyint(1) not null COMMENT 'Whether this access group allows manage',
    Foreign Key (venue_id) REFERENCES venues (venue_id),
    CONSTRAINT uq_vagid_venue_id UNIQUE (vagid, venue_id)
) COMMENT 'Venue Accesses Table';