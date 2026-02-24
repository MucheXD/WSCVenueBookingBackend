package applicationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func ReviewApplicationHandler(c *gin.Context) {
	applicationID, ok := parsePathInt(c, "application_id")
	if !ok {
		return
	}
	reviewerUID := c.GetString("UserID")
	vagid, sysPermMap, ok := getPermissionContext(c)
	if !ok {
		return
	}

	application, err := applicationSvc.GetApplicationByID(c.Request.Context(), applicationID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if !hasVenueApprovalPermission(vagid, sysPermMap, application.VenueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	var req reviewApplicationForm
	if err = c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	approval := models.ApplicationApproval{
		ApplicationID:  applicationID,
		Decision:       models.ApprovalDecision(req.Decision),
		KnownConflicts: req.KnownConflicts,
	}
	if req.Comment != nil {
		approval.Comment = &models.ApplicationComment{
			ID:          req.Comment.ID,
			CommentText: req.Comment.Text,
			Attachments: toAttachmentModelList(req.Comment.Attachments),
		}
	}

	result, err := applicationSvc.ReviewApplication(c.Request.Context(), approval, reviewerUID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	if len(result.NewConflicts) > 0 {
		utils.SetSuccessJsonResponse(c, map[string][]string{"new_conflicts": result.NewConflicts})
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}
