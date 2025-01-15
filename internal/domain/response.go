package domain

type APIBaseResponse struct {
	Code       int         `json:"code"`
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page      int16 `json:"page"`
	Limit     int8  `json:"limit"`
	TotalData int32 `json:"totalData"`
	TotalPage int16 `json:"totalPage"`
}

func NewResponse(code int, status bool, message string, data interface{}) APIBaseResponse {
	return APIBaseResponse{
		Code:       code,
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: nil,
	}
}

func (response APIBaseResponse) WithPagination(page int16, limit int8, total int16) APIBaseResponse {
	totalPage := (int(limit) + int(total) - 1) / int(limit)
	response.Pagination = &Pagination{
		Page:      page,
		Limit:     limit,
		TotalPage: int16(totalPage),
		TotalData: int32(total),
	}
	return response
}
