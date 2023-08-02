package params

import "time"

type Add struct {
	Page time.Time `json:"page" time_format:"2006-01-02" binding:"required" msg:"时间格式有问题"`
}
