package applicationSvc

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"gorm.io/gorm"
)

type ApprovalResult struct {
	NewConflicts []string
}

// cmt: 已迁移权限校验至 Controller 层，Service 入参移除了权限字段，仅处理申请业务与事务。
func CreateApplication(ctx context.Context, venueID int, applicantUID string, application models.Application) (int, error) {
	if len(application.TimeRequest) == 0 {
		return 0, ErrApplicationNoTimeRequest
	}
	for _, item := range application.TimeRequest {
		if !item.End.After(item.Begin) {
			return 0, ErrApplicationTimeRangeInvalid
		}
	}

	application.VenueID = venueID
	application.ApplicantUID = applicantUID
	application.ApplicationType = models.ApplicationTypeNormal
	application.ApplicationStatus = models.ApplicationStatusRequested

	var createdID int
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		id, err := repository.CreateApplicationWithTx(tx, &application)
		if err != nil {
			return err
		}
		createdID = id
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrApplicationCreateInDB, err)
	}
	return createdID, nil
}

// cmt: 已迁移权限校验至 Controller 层，Service 删除逻辑仅执行存在性确认与级联软删除。
func DeleteApplication(ctx context.Context, applicationID int) error {
	if _, err := repository.GetApplicationByID(applicationID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrApplicationNotFound
		}
		return fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repository.SoftDeleteApplicationCascadeWithTx(tx, applicationID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrApplicationDeleteInDB, err)
	}
	return nil
}

// cmt: 已迁移权限校验至 Controller 层，Service 列表接口仅按业务维度查询。
func ListVenueApplications(ctx context.Context, venueID int) ([]models.Application, error) {
	applications, err := repository.ListApplicationsByVenueID(venueID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	return applications, nil
}

func ListUserApplications(ctx context.Context, applicantUID string) ([]models.Application, error) {
	applications, err := repository.ListApplicationsByApplicantUID(applicantUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	return applications, nil
}

// cmt: 已迁移权限校验至 Controller 层，Service 审批逻辑仅校验状态与冲突并执行事务。
func ReviewApplication(ctx context.Context, approval models.ApplicationApproval, reviewerUID string) (ApprovalResult, error) {
	if approval.Decision != models.ApprovalDecisionApproved && approval.Decision != models.ApprovalDecisionRejected {
		return ApprovalResult{}, ErrApplicationDecisionInvalid
	}

	application, err := repository.GetApplicationByID(approval.ApplicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ApprovalResult{}, ErrApplicationNotFound
		}
		return ApprovalResult{}, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	if application.ApplicationStatus != models.ApplicationStatusRequested {
		return ApprovalResult{}, ErrApplicationStatusInvalid
	}

	comment := approval.Comment
	if comment != nil {
		comment.ApplicationID = approval.ApplicationID
		comment.CommenterUID = reviewerUID
	}

	if approval.Decision == models.ApprovalDecisionRejected {
		err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := repository.UpdateApplicationStatusWithTx(tx, approval.ApplicationID, models.ApplicationStatusRejected); err != nil {
				return err
			}
			if comment != nil {
				if _, err := repository.CreateApplicationCommentWithTx(tx, comment); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return ApprovalResult{}, fmt.Errorf("%w: %w", ErrApplicationUpdateInDB, err)
		}
		notifyApplicationDecision(application.ApplicantUID, approval.ApplicationID, models.ApplicationStatusRejected)
		return ApprovalResult{}, nil
	}

	conflictIDs, err := repository.GetConflictingApplicationIDs(approval.ApplicationID,
		[]string{models.ApplicationStatusRequested})
	if err != nil {
		return ApprovalResult{}, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}

	known := normalizeConflictIDs(approval.KnownConflicts)
	if !equalIntSlice(conflictIDs, known) {
		newConflicts := make([]string, 0, len(conflictIDs))
		for _, conflictID := range conflictIDs {
			newConflicts = append(newConflicts, strconv.Itoa(conflictID))
		}
		return ApprovalResult{NewConflicts: newConflicts}, nil
	}

	err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repository.UpdateApplicationStatusWithTx(tx, approval.ApplicationID, models.ApplicationStatusApproved); err != nil {
			return err
		}
		if err := repository.BatchRejectApplicationsWithTx(tx, conflictIDs); err != nil {
			return err
		}
		if comment != nil {
			if _, err := repository.CreateApplicationCommentWithTx(tx, comment); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return ApprovalResult{}, fmt.Errorf("%w: %w", ErrApplicationUpdateInDB, err)
	}

	notifyApplicationDecision(application.ApplicantUID, approval.ApplicationID, models.ApplicationStatusApproved)
	for _, conflictID := range conflictIDs {
		notifyApplicationDecision("", conflictID, models.ApplicationStatusRejected)
	}

	return ApprovalResult{}, nil
}

func GetApplicationByID(ctx context.Context, applicationID int) (*models.Application, error) {
	_ = ctx
	application, err := repository.GetApplicationByID(applicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	return application, nil
}

func normalizeConflictIDs(values []int) []int {
	if len(values) == 0 {
		return []int{}
	}
	set := map[int]struct{}{}
	result := make([]int, 0, len(values))
	for _, value := range values {
		if _, exists := set[value]; exists {
			continue
		}
		set[value] = struct{}{}
		result = append(result, value)
	}
	sort.Ints(result)
	return result
}

func equalIntSlice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func notifyApplicationDecision(applicantUID string, applicationID int, status string) {
	_ = applicantUID
	_ = applicationID
	_ = status
}
