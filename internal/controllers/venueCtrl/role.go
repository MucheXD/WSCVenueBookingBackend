package venueCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/roleSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type roleAccessForm struct {
	VenueID       int  `json:"venue_id" binding:"required"`
	AllowReserve  bool `json:"allow_reserve"`
	AllowApproval bool `json:"allow_approval"`
	AllowEdit     bool `json:"allow_edit"`
	AllowManage   bool `json:"allow_manage"`
}

type createRoleForm struct {
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Accesses    []roleAccessForm `json:"accesses"`
}

type updateRoleForm struct {
	Name        *string           `json:"name"`
	Description *string           `json:"description"`
	Accesses    *[]roleAccessForm `json:"accesses"`
}

type roleItemDTO struct {
	VAGID       int    `json:"vagid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListRolesHandler 列出场地权限角色组
// GET /api/role
func ListRolesHandler(c *gin.Context) {
	roles, err := roleSvc.ListRoles(c.Request.Context())
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	result := make([]roleItemDTO, 0, len(roles))
	for _, role := range roles {
		result = append(result, roleItemDTO{
			VAGID:       role.VAGID,
			Name:        role.Name,
			Description: role.Description,
		})
	}
	utils.SetSuccessJsonResponse(c, result)
}

// CreateRoleHandler 新建场地权限角色
// POST /api/role
func CreateRoleHandler(c *gin.Context) {
	var req createRoleForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	vagid, err := roleSvc.CreateRole(c.Request.Context(), roleSvc.CreateRoleInput{
		Name:        req.Name,
		Description: req.Description,
		Accesses:    toSvcAccesses(req.Accesses),
	})
	if err != nil {
		handleRoleSvcError(c, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]int{"vagid": vagid})
}

// UpdateRoleHandler 更改场地权限角色
// PUT /api/role/:vagid
func UpdateRoleHandler(c *gin.Context) {
	vagid, err := strconv.Atoi(c.Param("vagid"))
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var req updateRoleForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	input := roleSvc.UpdateRoleInput{
		Name:        req.Name,
		Description: req.Description,
	}
	if req.Accesses != nil {
		converted := toSvcAccesses(*req.Accesses)
		input.Accesses = &converted
	}

	if err := roleSvc.UpdateRole(c.Request.Context(), vagid, input); err != nil {
		handleRoleSvcError(c, err)
		return
	}

	utils.SetSuccessJsonResponse(c, nil)
}

func toSvcAccesses(inputs []roleAccessForm) []roleSvc.VenueAccessEntryInput {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]roleSvc.VenueAccessEntryInput, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, roleSvc.VenueAccessEntryInput{
			VenueID:       input.VenueID,
			AllowReserve:  input.AllowReserve,
			AllowApproval: input.AllowApproval,
			AllowEdit:     input.AllowEdit,
			AllowManage:   input.AllowManage,
		})
	}
	return result
}

func handleRoleSvcError(c *gin.Context, err error) {
	if errors.Is(err, roleSvc.ErrRoleNotFound) {
		apiException.AbortWithException(c, apiException.NotFound, err)
		return
	}
	if errors.Is(err, roleSvc.ErrRoleNameRequired) ||
		errors.Is(err, roleSvc.ErrRoleVenueIDInvalid) ||
		errors.Is(err, roleSvc.ErrRoleAccessGroupNotMatched) {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	apiException.AbortWithException(c, apiException.ServerError, err)
}
