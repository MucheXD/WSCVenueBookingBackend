package controllers

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// BuildingDTO 楼区DTO
type BuildingDTO struct {
	BuildingID       int    `json:"building_id"`
	BuildingName     string `json:"building_name"`
	LocationCampusID int    `json:"location_campus_id"`
}

// CampusDTO 校区DTO
type CampusDTO struct {
	CampusID   int    `json:"campus_id"`
	CampusName string `json:"campus_name"`
}

// LocationInfo 位置信息（楼区和校区）
type LocationInfo struct {
	Buildings []BuildingDTO `json:"buildings"`
	Campuses  []CampusDTO   `json:"campuses"`
}

// GetVenueLocationsHandler 获取场地位置信息（楼区和校区）
// GET /api/venue/locations
func GetVenueLocationsHandler(c *gin.Context) {
	// 从上下文获取用户的VAGID（由AuthMiddleware设置）
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

	// 调用服务层获取楼区和校区信息
	buildings, campuses, err := venueSvc.GetAccessibleBuildingsAndCampuses(c.Request.Context(), vagid)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 转换为DTO
	buildingDTOs := make([]BuildingDTO, 0, len(buildings))
	for _, b := range buildings {
		buildingDTOs = append(buildingDTOs, BuildingDTO{
			BuildingID:       b.ID,
			BuildingName:     b.Name,
			LocationCampusID: b.CampusID,
		})
	}

	campusDTOs := make([]CampusDTO, 0, len(campuses))
	for _, c := range campuses {
		campusDTOs = append(campusDTOs, CampusDTO{
			CampusID:   c.ID,
			CampusName: c.Name,
		})
	}

	// 组装返回包体
	locationInfo := LocationInfo{
		Buildings: buildingDTOs,
		Campuses:  campusDTOs,
	}

	utils.SetSuccessJsonResponse(c, locationInfo)
}
