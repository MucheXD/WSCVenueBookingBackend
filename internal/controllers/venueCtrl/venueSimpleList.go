package venueCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type VenueSimpleDTO struct {
	VenueID     int      `json:"venue_id"`
	Name        string   `json:"name"`
	BuildingID  int      `json:"building_id"`
	TypeID      int      `json:"type_id"`
	Permissions []string `json:"permissions"`
}

// ListVenueAccessBodiesHandler 列出可修改权限的场地（轻量数据）
// GET /api/venue/list
func ListVenueAccessBodiesHandler(c *gin.Context) {
	vagidVal, exists := c.Get("VenueAccessGroupID")
	if !exists {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}
	vagid, ok := vagidVal.(int)
	if !ok {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	venues, err := venueSvc.ListVenueBodies(c.Request.Context())
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	result := make([]VenueSimpleDTO, 0, len(venues))
	for _, venue := range venues {
		result = append(result, VenueSimpleDTO{
			VenueID:     venue.ID,
			Name:        venue.Name,
			BuildingID:  venue.BuildingID,
			TypeID:      venue.TypeID,
			Permissions: getVenuePermissionStrings(vagid, venue.ID),
		})
	}

	utils.SetSuccessJsonResponse(c, result)
}
