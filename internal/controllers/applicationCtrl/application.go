package applicationCtrl

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"github.com/gin-gonic/gin"
)

// 传入申请单 -> 申请单模型
func toApplicationModel(req createApplicationForm) (models.Application, error) {
	timeslots := make([]models.ApplicationTimeRequest, 0, len(req.TimeRequest))
	for _, item := range req.TimeRequest { // 解析 Timeslots JSONArray 为模型
		begin, err := time.Parse(time.RFC3339, item.Begin)
		if err != nil {
			return models.Application{}, err
		}
		end, err := time.Parse(time.RFC3339, item.End)
		if err != nil {
			return models.Application{}, err
		}
		timeslots = append(timeslots, models.ApplicationTimeRequest{Begin: begin, End: end})
	}

	return models.Application{
		ApplicationType:        models.ApplicationTypeNormal,
		TimeRequest:            timeslots,
		EstimatedParticipants:  req.EstimatedParticipants,
		DescriptionText:        req.Description,
		Attachments:            toAttachmentModelList(req.Attachments),
		ActivityName:           req.ActivityName,
		ActivityOrganizer:      req.ActivityOrganizer,
		ActivityCoordinatorRaw: marshalCoordinator(req.ActivityCoordinator),
	}, nil
}

// 申请单模型 -> 申请单传出
func toApplicationResponseList(applications []models.Application) []applicationResponseDTO {
	result := make([]applicationResponseDTO, 0, len(applications))
	for _, application := range applications {
		timeslots := make([]timetableDTO, 0, len(application.TimeRequest))
		for _, slot := range application.TimeRequest {
			timeslots = append(timeslots, timetableDTO{
				Begin: slot.Begin.Format(time.RFC3339),
				End:   slot.End.Format(time.RFC3339),
			})
		}

		comments := make([]commentDTO, 0, len(application.Comments))
		for _, comment := range application.Comments {
			comments = append(comments, commentDTO{
				ID:          comment.ID,
				Text:        comment.CommentText,
				Attachments: toAttachmentDTOList(comment.Attachments),
			})
		}

		result = append(result, applicationResponseDTO{
			ID:                    application.ID,
			ApplicationType:       models.ApplicationTypeNormal,
			ApplicationStatus:     application.ApplicationStatus,
			TimeRequest:           timeslots,
			EstimatedParticipants: application.EstimatedParticipants,
			Description:           application.DescriptionText,
			Attachments:           toAttachmentDTOList(application.Attachments),
			ActivityName:          application.ActivityName,
			ActivityOrganizer:     application.ActivityOrganizer,
			ActivityCoordinator:   unmarshalCoordinator(application.ActivityCoordinatorRaw),
			Comments:              comments,
		})
	}
	return result
}

// 传入附件表 -> 附件表模型
func toAttachmentModelList(values []attachmentDTO) []models.Attachment {
	if len(values) == 0 {
		return []models.Attachment{}
	}
	result := make([]models.Attachment, 0, len(values))
	for _, value := range values {
		result = append(result, models.Attachment{
			Index:       value.Index,
			FileToken:   value.FileToken,
			BizFileType: value.FileType,
			BizFileName: value.FileName,
		})
	}
	return result
}

// 附件表模型 -> 附件表传出
func toAttachmentDTOList(values []models.Attachment) []attachmentDTO {
	if len(values) == 0 {
		return []attachmentDTO{}
	}
	result := make([]attachmentDTO, 0, len(values))
	for _, value := range values {
		result = append(result, attachmentDTO{
			Index:     value.Index,
			FileToken: value.FileToken,
			FileType:  value.BizFileType,
			FileName:  value.BizFileName,
		})
	}
	return result
}

// 辅助函数，从路径参数中解析整数，自动处理错误
func parsePathInt(c *gin.Context, key string) (int, bool) {
	value, err := strconv.Atoi(c.Param(key))
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return 0, false
	}
	return value, true
}

// 辅助函数，判断用户权限字段是否存在于上下文，并返回对应的值
func getPermissionContext(c *gin.Context) (int, uint64, bool) {
	vagidVal, exists := c.Get("VenueAccessGroupID")
	if !exists {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	vagid, ok := vagidVal.(int)
	if !ok {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}

	sysPermVal, exists := c.Get("SysPermissionMap")
	if !exists {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	sysPermMap, ok := sysPermVal.(uint64)
	if !ok {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	return vagid, sysPermMap, true
}

// 统一负责转换 Service 层错误与 apiException
func handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, applicationSvc.ErrApplicationNotFound) {
		apiException.AbortWithException(c, apiException.NotFound, err)
		return
	}
	if errors.Is(err, applicationSvc.ErrApplicationNoTimeRequest) ||
		errors.Is(err, applicationSvc.ErrApplicationTimeRangeInvalid) ||
		errors.Is(err, applicationSvc.ErrApplicationDecisionInvalid) ||
		errors.Is(err, applicationSvc.ErrApplicationStatusInvalid) {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	apiException.AbortWithException(c, apiException.ServerError, err)
}

// cmt: 额外重构，新增 Controller 层权限函数，统一处理申请模块权限判断逻辑。
func hasVenueReservePermission(vagid int, sysPermMap uint64, venueID int) bool {
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueReservation) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Reserve)
}

func hasVenueApprovalPermission(vagid int, sysPermMap uint64, venueID int) bool {
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueApproval) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Approval)
}

func canDeleteApplication(requesterUID string, vagid int, sysPermMap uint64, application models.Application) bool {
	if requesterUID == application.ApplicantUID {
		return true
	}
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueManage) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, application.VenueID, venuePermission.Manage)
}

// 预留处理函数
func marshalCoordinator(values []any) []byte {
	if len(values) == 0 {
		return nil
	}
	b, err := json.Marshal(values)
	if err != nil {
		return nil
	}
	return b
}

// 预留处理函数
func unmarshalCoordinator(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	var values any
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil
	}
	return values
}
