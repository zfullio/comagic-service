package entity

type DynamicCallTracking struct {
	ReservationTime     string  `json:"reservation_time"`
	CountVirtualNumbers int64   `json:"count_virtual_numbers"`
	CountVisits         int64   `json:"count_visits"`
	CoverageVisitors    float64 `json:"coverage_visitors"`
}
