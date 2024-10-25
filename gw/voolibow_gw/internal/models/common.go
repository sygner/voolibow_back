package models

type Pagination struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
	Total  bool  `json:"get_total"`
}
