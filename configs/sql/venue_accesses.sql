DROP TABLE IF EXISTS venue_accesses;

CREATE TABLE venue_accesses (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'Unique identifier for each access info, NOT for use',
    vagid int not null COMMENT 'Unique identifier for each access group, 0 for user default value',
    venue_id int not null COMMENT 'Target venue id of this access info',
    allow_reserve BOOLEAN not null DEFAULT false COMMENT 'Whether this access group allows reservation',
    allow_approval BOOLEAN not null DEFAULT false COMMENT 'Whether this access group allows approval',
    # Edit 权限包含修改场地名称与描述、修改场地图片、场地停用、场地移除等
    allow_edit BOOLEAN not null DEFAULT false COMMENT 'Whether this access group allows edit',
    # Manage 权限包含设置维护时段、维护可用设备列表、发布场地公告等
    allow_manage BOOLEAN not null DEFAULT false COMMENT 'Whether this access group allows manage',
    Foreign Key (venue_id) REFERENCES venues (venue_id),
    Foreign Key (vagid) REFERENCES venue_roles (vagid),
    CONSTRAINT uq_vagid_venue_id UNIQUE (vagid, venue_id)
) COMMENT 'Venue Accesses Table';