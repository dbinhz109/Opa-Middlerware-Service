package model

type ListParam struct {
	TimeEnd     *int64                  `json:"time_end,omitempty"`
	TimeStart   *int64                  `json:"time_start,omitempty"`
	Offset      int64                   `json:"offset,omitempty"`
	Limit       int64                   `json:"limit,omitempty"`
	Count       bool                    `json:"count,omitempty"`
	Sort        *string                 `json:"sort,omitempty"`
	Size        int64                   `json:"size,omitempty"`
	SortAf      []interface{}           `json:"sort_af,omitempty"`
	SearchAfter []any                   `json:"search_after,omitempty"`
	Aggs        *map[string]interface{} `json:"aggs,omitempty"`
	AggsType    *string                 `json:"aggs_type,omitempty"`
}

func Ptr[T any](v T) *T {
	return &v
}
