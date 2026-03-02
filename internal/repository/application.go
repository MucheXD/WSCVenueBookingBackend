package repository

import (
	"sort"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type ApplicationEntity struct {
	ID                    int            `gorm:"column:id;primaryKey"`
	VenueID               int            `gorm:"column:venue_id"`
	ApplicantUID          string         `gorm:"column:applicant_uid"`
	ApplicationStatus     string         `gorm:"column:application_status"`
	DescriptionText       string         `gorm:"column:description_text"`
	EstimatedParticipants *int           `gorm:"column:estimated_participants"`
	HasAttachments        bool           `gorm:"column:has_attachments"`
	ActivityName          string         `gorm:"column:activity_name"`
	ActivityOrganizer     string         `gorm:"column:activity_organizer"`
	ActivityCoordinator   []byte         `gorm:"column:activity_coordinator"`
	CreatedAt             time.Time      `gorm:"column:created_at;autoCreateTime"`
	ApprovalAt            *time.Time     `gorm:"column:approval_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ApplicationEntity) TableName() string {
	return "applications"
}

type ApplicationCommentEntity struct {
	ID             int            `gorm:"column:id;primaryKey"`
	ApplicationID  int            `gorm:"column:application_id"`
	CommenterUID   string         `gorm:"column:commenter_uid"`
	CommentText    string         `gorm:"column:comment_text"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	HasAttachments bool           `gorm:"column:has_attachments"`
}

func (ApplicationCommentEntity) TableName() string {
	return "application_comments"
}

type ApplicationTimeslotEntity struct {
	ID            int            `gorm:"column:id;primaryKey"`
	VenueID       int            `gorm:"column:venue_id"`
	StartTime     time.Time      `gorm:"column:start_time"`
	EndTime       *time.Time     `gorm:"column:end_time"`
	ApplicationID *int           `gorm:"column:application_id"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ApplicationTimeslotEntity) TableName() string {
	return "venue_timeslots"
}

func CreateApplicationTx(tx *gorm.DB, application *models.Application) (int, error) {
	entity := ApplicationEntity{
		VenueID:               application.VenueID,
		ApplicantUID:          application.ApplicantUID,
		ApplicationStatus:     application.ApplicationStatus,
		DescriptionText:       application.DescriptionText,
		EstimatedParticipants: application.EstimatedParticipants,
		HasAttachments:        len(application.Attachments) > 0,
		ActivityName:          application.ActivityName,
		ActivityOrganizer:     application.ActivityOrganizer,
		ActivityCoordinator:   application.ActivityCoordinatorRaw,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}

	if len(application.TimeRequest) > 0 {
		timeslotEntities := make([]ApplicationTimeslotEntity, 0, len(application.TimeRequest))
		for _, slot := range application.TimeRequest {
			timeslotEntities = append(timeslotEntities, ApplicationTimeslotEntity{
				VenueID:       application.VenueID,
				StartTime:     slot.Begin,
				EndTime:       &slot.End,
				ApplicationID: &entity.ID,
			})
		}
		if err := tx.Create(&timeslotEntities).Error; err != nil {
			return 0, err
		}
	}

	if len(application.Attachments) > 0 {
		attachments := make([]models.Attachment, 0, len(application.Attachments))
		for idx, attachment := range application.Attachments {
			attachment.Index = idx
			attachments = append(attachments, attachment)
		}
		if err := CreateAttachmentsTx(tx, AttachmentBizTypeApplication, entity.ID, attachments); err != nil {
			return 0, err
		}
	}

	return entity.ID, nil
}

func CreateApplicationCommentTx(tx *gorm.DB, comment *models.ApplicationComment) (int, error) {
	entity := ApplicationCommentEntity{
		ApplicationID:  comment.ApplicationID,
		CommenterUID:   comment.CommenterUID,
		CommentText:    comment.CommentText,
		HasAttachments: len(comment.Attachments) > 0,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}
	if len(comment.Attachments) > 0 {
		attachments := make([]models.Attachment, 0, len(comment.Attachments))
		for idx, attachment := range comment.Attachments {
			attachment.Index = idx
			attachments = append(attachments, attachment)
		}
		if err := CreateAttachmentsTx(tx, AttachmentBizTypeApplicationComment, entity.ID, attachments); err != nil {
			return 0, err
		}
	}
	return entity.ID, nil
}

func GetFullApplicationByID(applicationID int) (*models.Application, error) {
	applications, err := queryApplications(database.DB.Where(&ApplicationEntity{ID: applicationID}))
	if err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &applications[0], nil
}

func GetApplicationBodyByID(applicationID int) (*models.Application, error) {
	var entity ApplicationEntity
	err := database.DB.
		Model(&ApplicationEntity{}).
		Where("id = ?", applicationID).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &models.Application{
		ID:                     entity.ID,
		VenueID:                entity.VenueID,
		ApplicantUID:           entity.ApplicantUID,
		ApplicationStatus:      entity.ApplicationStatus,
		TimeRequest:            []models.ApplicationTimeRequest{},
		EstimatedParticipants:  entity.EstimatedParticipants,
		DescriptionText:        entity.DescriptionText,
		ActivityName:           entity.ActivityName,
		ActivityOrganizer:      entity.ActivityOrganizer,
		HasAttachments:         entity.HasAttachments,
		Attachments:            []models.Attachment{},
		Comments:               []models.ApplicationComment{},
		ActivityCoordinatorRaw: entity.ActivityCoordinator,
		CreatedAt:              entity.CreatedAt,
		ApprovalAt:             entity.ApprovalAt,
	}, nil
}

func ListApplicationsByVenueID(venueID int) ([]models.Application, error) {
	return queryApplications(database.DB.Where(&ApplicationEntity{VenueID: venueID}))
}

func ListApplicationsByApplicantUID(applicantUID string) ([]models.Application, error) {
	return queryApplications(database.DB.Where(&ApplicationEntity{ApplicantUID: applicantUID}))
}

// 软删除申请单及其关联的时间段、评论、附件等数据
func SoftDeleteApplicationsTx(tx *gorm.DB, applicationID int) error {
	// 删除申请单关联时间段
	if err := tx.Model(&ApplicationTimeslotEntity{}).
		Where(&ApplicationTimeslotEntity{ApplicationID: &applicationID}).
		Delete(&ApplicationTimeslotEntity{}).Error; err != nil {
		return err
	}
	// 删除申请单附件
	if err := SoftDeleteBizAttachmentsTx(tx, AttachmentBizTypeApplication, []int{applicationID}); err != nil {
		return err
	}

	// 删除关联评论及评论附件
	var comments []ApplicationCommentEntity
	// 列出评论
	if err := tx.Model(&ApplicationCommentEntity{}).
		Where(&ApplicationCommentEntity{ApplicationID: applicationID}).
		Find(&comments).Error; err != nil {
		return err
	}
	// 删除评论附件
	commentIDs := make([]int, 0, len(comments))
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}

	if err := SoftDeleteBizAttachmentsTx(tx, AttachmentBizTypeApplicationComment, commentIDs); err != nil {
		return err
	}
	// 删除评论主体
	if err := tx.Model(&ApplicationCommentEntity{}).
		Where(&ApplicationCommentEntity{ApplicationID: applicationID}).
		Delete(&ApplicationCommentEntity{}).Error; err != nil {
		return err
	}

	// 删除申请单主体
	return tx.Model(&ApplicationEntity{}).
		Where(&ApplicationEntity{ID: applicationID}).
		Delete(&ApplicationEntity{}).Error
}

func UpdateApplicationStatusTx(tx *gorm.DB, applicationID int, status string) error {
	return UpdateApplicationStatusesTx(tx, []int{applicationID}, status)
}

// UpdateApplicationStatusesTx 批量更新申请单状态
// currentStatuses 可选：用于限制仅更新当前状态在该集合内的记录
// 若新状态为已批准或已拒绝，同步将 approval_at 设置为当前时间
func UpdateApplicationStatusesTx(tx *gorm.DB, applicationIDs []int, status string, currentStatuses ...string) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	query := tx.Model(&ApplicationEntity{}).Where("id IN ?", applicationIDs)
	// 动态注入额外条件
	if len(currentStatuses) > 0 {
		query = query.Where("application_status IN ?", currentStatuses)
	}
	updates := map[string]interface{}{
		"application_status": status,
	}
	if status == models.ApplicationStatusApproved || status == models.ApplicationStatusRejected {
		now := time.Now()
		updates["approval_at"] = now
	}
	return query.Updates(updates).Error
}

// 获取冲突的申请单ID，使用高级数据库查询完成
// candidateStatuses 为从申请单状态过滤器，例如：仅将“req”状态的从申请单定义为冲突
func GetConflictingApplicationIDs(applicationID int, candidateStatuses []string) ([]int, error) {
	if len(candidateStatuses) == 0 {
		return []int{}, nil
	}

	var conflictIDs []int
	err := database.DB.
		Table("applications a1").                                              // 操作表：applications 定义 a1 为主申请单
		Where("a1.id = ?", applicationID).                                     // 筛选：主申请单ID对应传入ID
		Distinct("a2.id").                                                     // 置唯一：相同从申请单以ID为依据仅输出一次
		Joins("JOIN venue_timeslots t1 ON t1.application_id = a1.id").         // 提取：t1 为主申请单时间段
		Joins("JOIN venue_timeslots t2 ON t2.venue_id = a1.venue_id").         // 提取：t2 为与主申请单申请场地相同的时间段
		Joins("JOIN applications a2 ON a2.id = t2.application_id").            // 提取：a2 为 t2 关连的申请单(即从申请单)
		Where("a2.id <> a1.id").                                               // 排除：从申请单不包括主申请单
		Where("a1.deleted_at IS NULL").                                        // 复杂查询需手动排除软删除
		Where("a2.deleted_at IS NULL").                                        // 复杂查询需手动排除软删除
		Where("t1.deleted_at IS NULL").                                        // 复杂查询需手动排除软删除
		Where("t2.deleted_at IS NULL").                                        // 复杂查询需手动排除软删除
		Where("a2.application_status IN ?", candidateStatuses).                // 筛选：具有目标状态的从申请单
		Where("t1.start_time < COALESCE(t2.end_time, '9999-12-31 23:59:59')"). // 筛选：时间区段左重叠判断，为可能空的 end_time 定义默认值
		Where("COALESCE(t1.end_time, '9999-12-31 23:59:59') > t2.start_time"). // 筛选：时间区段右重叠判断，为可能空的 end_time 定义默认值
		Pluck("a2.id", &conflictIDs).Error                                     // 输出：符合上述条件从申请单ID
	if err != nil {
		return nil, err
	}
	sort.Ints(conflictIDs)
	return conflictIDs, nil
}

// 以传入 scope 的条件查询所有申请单，同步查询时间表、评论表、附件表，并进行组装
// 附加条件时请传入对应条件的 scope
func queryApplications(scope *gorm.DB) ([]models.Application, error) {
	var appEntities []ApplicationEntity

	// 查询申请单主体
	if err := scope.
		Model(&ApplicationEntity{}).
		Order("id DESC").
		Find(&appEntities).Error; err != nil {
		return nil, err
	}
	if len(appEntities) == 0 {
		return []models.Application{}, nil
	}

	// 获取所有申请单ID以便后续查询关联数据
	appIDs := make([]int, 0, len(appEntities))
	for _, app := range appEntities {
		appIDs = append(appIDs, app.ID)
	}

	// 批量查询申请单对应时间表
	var timeslotEntities []ApplicationTimeslotEntity
	if err := database.DB.
		Model(&ApplicationTimeslotEntity{}).
		Where("application_id IN ?", appIDs).
		Order("application_id ASC, start_time ASC").
		Find(&timeslotEntities).Error; err != nil {
		return nil, err
	}
	timeslotMap := make(map[int][]models.ApplicationTimeRequest)
	for _, timeslot := range timeslotEntities {
		// ENH: 此处过滤了未确定结束的时间段，后续可考虑增加特殊处理逻辑以供前端特殊显示
		if timeslot.ApplicationID == nil || timeslot.EndTime == nil {
			continue
		}
		appID := *timeslot.ApplicationID
		timeslotMap[appID] = append(timeslotMap[appID], models.ApplicationTimeRequest{
			Begin: timeslot.StartTime,
			End:   *timeslot.EndTime,
		})
	}

	// 批量查询申请单对应附件
	var appAttachmentEntities []AttachmentEntity
	if err := database.DB.
		Model(&AttachmentEntity{}).
		Where("biz_type = ?", AttachmentBizTypeApplication).
		Where("biz_id IN ?", appIDs).
		Order("biz_id ASC, biz_index ASC").
		Find(&appAttachmentEntities).Error; err != nil {
		return nil, err
	}
	appAttachmentMap := make(map[int][]models.Attachment)
	for _, attachment := range appAttachmentEntities {
		appAttachmentMap[attachment.BizID] = append(appAttachmentMap[attachment.BizID], attachment.toDomain())
	}

	// 批量查询申请单对应评论
	var commentEntities []ApplicationCommentEntity
	if err := database.DB.
		Model(&ApplicationCommentEntity{}).
		Where("application_id IN ?", appIDs).
		Order("application_id ASC, created_at ASC").
		Find(&commentEntities).Error; err != nil {
		return nil, err
	}
	commentIDs := make([]int, 0, len(commentEntities))
	for _, comment := range commentEntities { // 获取评论ID以便后续查询评论附件
		commentIDs = append(commentIDs, comment.ID)
	}
	commentAttachmentMap := make(map[int][]models.Attachment)
	if len(commentIDs) > 0 {
		var commentAttachmentEntities []AttachmentEntity
		if err := database.DB.
			Model(&AttachmentEntity{}).
			Where("biz_type = ?", AttachmentBizTypeApplicationComment).
			Where("biz_id IN ?", commentIDs).
			Order("biz_id ASC, biz_index ASC").
			Find(&commentAttachmentEntities).Error; err != nil {
			return nil, err
		}
		for _, attachment := range commentAttachmentEntities {
			commentAttachmentMap[attachment.BizID] = append(commentAttachmentMap[attachment.BizID], attachment.toDomain())
		}
	}
	commentsByApplication := make(map[int][]models.ApplicationComment)
	for _, comment := range commentEntities { // 组装评论和评论附件为模型
		commentsByApplication[comment.ApplicationID] = append(commentsByApplication[comment.ApplicationID], models.ApplicationComment{
			ID:             comment.ID,
			ApplicationID:  comment.ApplicationID,
			CommenterUID:   comment.CommenterUID,
			CommentText:    comment.CommentText,
			CommentTime:    comment.CreatedAt,
			HasAttachments: comment.HasAttachments,
			Attachments:    commentAttachmentMap[comment.ID],
		})
	}

	// 最终组装申请单模型
	applications := make([]models.Application, 0, len(appEntities))
	for _, entity := range appEntities {
		applicationType := models.ApplicationTypeNormal
		applications = append(applications, models.Application{
			ID:                     entity.ID,
			VenueID:                entity.VenueID,
			ApplicantUID:           entity.ApplicantUID,
			ApplicationType:        applicationType,
			ApplicationStatus:      entity.ApplicationStatus,
			TimeRequest:            timeslotMap[entity.ID],
			EstimatedParticipants:  entity.EstimatedParticipants,
			DescriptionText:        entity.DescriptionText,
			ActivityName:           entity.ActivityName,
			ActivityOrganizer:      entity.ActivityOrganizer,
			HasAttachments:         entity.HasAttachments,
			Attachments:            appAttachmentMap[entity.ID],
			Comments:               commentsByApplication[entity.ID],
			ActivityCoordinatorRaw: entity.ActivityCoordinator,
			CreatedAt:              entity.CreatedAt,
			ApprovalAt:             entity.ApprovalAt,
		})
	}
	return applications, nil
}

// statsRow 用于接收单行聚合查询结果
type statsRow struct {
	AllApplications   int64 `gorm:"column:all_applications"`
	AllApproved       int64 `gorm:"column:all_approved"`
	AllRejected       int64 `gorm:"column:all_rejected"`
	AllPending        int64 `gorm:"column:all_pending"`
	Last7Applications int64 `gorm:"column:last7_applications"`
	Last7Approved     int64 `gorm:"column:last7_approved"`
	Last7Rejected     int64 `gorm:"column:last7_rejected"`
}

// GetApplicationStats 通过单次聚合查询获取申请单全量与近7天统计数据
func GetApplicationStats() (*models.ApplicationStats, error) {
	cutoff := time.Now().AddDate(0, 0, -7)

	var row statsRow
	err := database.DB.
		Model(&ApplicationEntity{}).
		Select(`
			COUNT(*) AS all_applications,
			SUM(CASE WHEN application_status = ? THEN 1 ELSE 0 END) AS all_approved,
			SUM(CASE WHEN application_status = ? THEN 1 ELSE 0 END) AS all_rejected,
			SUM(CASE WHEN application_status = ? THEN 1 ELSE 0 END) AS all_pending,
			SUM(CASE WHEN created_at >= ? THEN 1 ELSE 0 END) AS last7_applications,
			SUM(CASE WHEN application_status = ? AND approval_at >= ? THEN 1 ELSE 0 END) AS last7_approved,
			SUM(CASE WHEN application_status = ? AND approval_at >= ? THEN 1 ELSE 0 END) AS last7_rejected
		`,
			models.ApplicationStatusApproved,
			models.ApplicationStatusRejected,
			models.ApplicationStatusRequested,
			cutoff,
			models.ApplicationStatusApproved, cutoff,
			models.ApplicationStatusRejected, cutoff,
		).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &models.ApplicationStats{
		AllApplications:   row.AllApplications,
		AllApproved:       row.AllApproved,
		AllRejected:       row.AllRejected,
		AllPending:        row.AllPending,
		Last7Applications: row.Last7Applications,
		Last7Approved:     row.Last7Approved,
		Last7Rejected:     row.Last7Rejected,
	}, nil
}
