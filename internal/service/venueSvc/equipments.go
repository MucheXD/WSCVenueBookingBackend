package venueSvc

import (
	"encoding/json"
	"strings"
)

type venueEquipmentInput struct {
	Name      *string `json:"n"`
	Amount    *int    `json:"a,omitempty"`
	Available *bool   `json:"v,omitempty"`
	Comment   *string `json:"c,omitempty"`
}

type venueEquipmentNormalized struct {
	Name      string  `json:"n"`
	Amount    int     `json:"a"`
	Available bool    `json:"v"`
	Comment   *string `json:"c,omitempty"`
}

func normalizeVenueEquipments(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var input []venueEquipmentInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, ErrVenueEquipmentsInvalid
	}

	/*
		校验 JSON 格式合法性：主体应为一个 Array，存储一系列 Object，其中每个 Object 应该为：
		{
		    "n": "设备前端显示名称... 此项必有",
		    "a": (int)该种设备数量... 此项库内必有，但传入可缺省，若缺省则给默认值1,
		    "v": (bool)设备可用性... 此项库内必有，但传入可缺省，若缺省则则给默认值true,
		    "c": "注释... 此项可无，用于前端鼠标悬停显示"
		}
	*/

	normalized := make([]venueEquipmentNormalized, 0, len(input))
	for _, item := range input {

		if item.Name == nil || strings.TrimSpace(*item.Name) == "" {
			return nil, ErrVenueEquipmentsInvalid // 设备名称：不可空
		}

		amount := 1 // 默认数量为 1
		if item.Amount != nil {
			amount = *item.Amount
		}

		available := true // 默认可用性为 true
		if item.Available != nil {
			available = *item.Available
		}

		normalized = append(normalized, venueEquipmentNormalized{
			Name:      strings.TrimSpace(*item.Name),
			Amount:    amount,
			Available: available,
			Comment:   item.Comment,
		})
	}

	b, err := json.Marshal(normalized) // 编译回 JSON 存储
	if err != nil {
		return nil, ErrVenueEquipmentsInvalid
	}
	return b, nil
}
