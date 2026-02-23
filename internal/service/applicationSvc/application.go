package applicationSvc

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"gorm.io/gorm"
)

type ApprovalResult struct {
	NewConflicts []string
}

func CreateApplication(ctx context.Context, venueID int, applicantUID string, vagid int, sysPermMap uint64, application models.Application) (int, error) {
	if !canReserve(venueID, vagid, sysPermMap) {
		return 0, ErrApplicationPermissionDenied
	}
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

func DeleteApplication(ctx context.Context, applicationID int, requesterUID string, vagid int, sysPermMap uint64) error {
	application, err := repository.GetApplicationByID(applicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrApplicationNotFound
		}
		return fmt.Errorf("%w: %w", ErrApplicationQueryInDB, err)
	}

	if !canDeleteApplication(*application, requesterUID, vagid, sysPermMap) {
		return ErrApplicationPermissionDenied
	}

	err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func ListVenueApplications(ctx context.Context, venueID int, vagid int, sysPermMap uint64) ([]models.Application, error) {
	if !canReserve(venueID, vagid, sysPermMap) {
		return nil, ErrApplicationPermissionDenied
	}
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

func ReviewApplication(ctx context.Context, approval models.ApplicationApproval, reviewerUID string, vagid int, sysPermMap uint64) (ApprovalResult, error) {
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
	if !canApprove(application.VenueID, vagid, sysPermMap) {
		return ApprovalResult{}, ErrApplicationPermissionDenied
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

func canReserve(venueID int, vagid int, sysPermMap uint64) bool {
	// TODO 此函数可以被提取
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueReservation) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Reserve)
}

func canApprove(venueID int, vagid int, sysPermMap uint64) bool {
	// TODO 此函数可以被提取
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueApproval) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Approval)
}

func canDeleteApplication(application models.Application, requesterUID string, vagid int, sysPermMap uint64) bool {
	if requesterUID == application.ApplicantUID {
		return true
	}
	if systemPermission.Check(sysPermMap, systemPermission.AllVenueManage) {
		return true
	}
	return venuePermission.CheckVenuePermission(vagid, application.VenueID, venuePermission.Manage)
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
