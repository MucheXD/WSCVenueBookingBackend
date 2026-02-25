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
		id, err := repository.CreateApplicationTx(tx, &application)
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

func DeleteApplication(ctx context.Context, applicationID int) error {
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repository.SoftDeleteApplicationsTx(tx, applicationID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrApplicationNotFound
		}
		return fmt.Errorf("%w: %w", ErrApplicationDeleteInDB, err)
	}
	return nil
}

func GetFullApplicationByID(ctx context.Context, applicationID int) (*models.Application, error) {
	_ = ctx
	application, err := repository.GetFullApplicationByID(applicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	return application, nil
}

func GetApplicationBodyByID(ctx context.Context, applicationID int) (*models.Application, error) {
	_ = ctx
	application, err := repository.GetApplicationBodyByID(applicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}
	return application, nil
}

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

func ReviewApplication(ctx context.Context, approval models.ApplicationApproval, reviewerUID string) (ApprovalResult, error) {

	// 参数合法性检查，仅允许批准或拒绝两种决策
	if approval.Decision != models.ApplicationStatusApproved && approval.Decision != models.ApplicationStatusRejected {
		return ApprovalResult{}, ErrApplicationDecisionInvalid
	}

	application, err := repository.GetApplicationBodyByID(approval.ApplicationID)
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
		comment.ID = 0 // 确保 Comment ID 在创建时由数据库自增生成
		comment.ApplicationID = approval.ApplicationID
		comment.CommenterUID = reviewerUID
	}

	// 处理拒绝申请单的情况：直接更新申请单状态并记录审批意见，无需检查冲突
	if approval.Decision == models.ApplicationStatusRejected {
		err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := repository.UpdateApplicationStatusTx(tx, approval.ApplicationID, models.ApplicationStatusRejected); err != nil {
				return err
			}
			if comment != nil {
				if _, err := repository.CreateApplicationCommentTx(tx, comment); err != nil {
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

	// 处理批准申请单的情况：首先检查是否存在新的冲突，如果存在则返回新冲突列表；如果不存在则更新申请单状态、自动拒绝冲突的申请单、记录审批意见
	conflictIDs, err := repository.GetConflictingApplicationIDs(approval.ApplicationID,
		[]string{models.ApplicationStatusRequested})
	if err != nil {
		return ApprovalResult{}, fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}

	// 对冲突ID列表进行预处理（去重与排序），因为 equalIntSlice 是单对单检查
	conflictIDs = normalizeConflictIDs(conflictIDs)
	known := normalizeConflictIDs(approval.KnownConflicts)

	if !equalIntSlice(conflictIDs, known) {
		newConflicts := make([]string, 0, len(conflictIDs))
		for _, conflictID := range conflictIDs {
			newConflicts = append(newConflicts, strconv.Itoa(conflictID))
		}
		// 若冲突列表与已知冲突列表不一致，则返回新冲突列表，前端可据此提示用户确认是否继续审批
		return ApprovalResult{NewConflicts: newConflicts}, nil
	}

	// 开始数据库事务更新：更改申请单状态、自动拒绝冲突的申请单、记录审批意见
	err = database.DB.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := repository.UpdateApplicationStatusTx(tx, approval.ApplicationID, models.ApplicationStatusApproved); err != nil {
				return err
			}
			if err := repository.UpdateApplicationStatusesTx(
				tx,
				conflictIDs,
				models.ApplicationStatusRejected,
				models.ApplicationStatusRequested,
			); err != nil {
				return err
			}
			if comment != nil {
				if _, err := repository.CreateApplicationCommentTx(tx, comment); err != nil {
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

// 对输入的原始冲突ID列表进行“去重”和“排序”处理
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

// 预留函数，等待“站内信”功能实现
func notifyApplicationDecision(applicantUID string, applicationID int, status string) {
	_ = applicantUID
	_ = applicationID
	_ = status
}
