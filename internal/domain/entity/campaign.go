package entity

import "time"

type Campaign struct {
	ID                  int64               `json:"id"`
	Status              string              `json:"status"`
	CreationTime        string              `json:"creation_time"`
	Description         string              `json:"description"`
	SiteID              int64               `json:"site_id"`
	SiteDomainName      string              `json:"site_domain_name"`
	Costs               float64             `json:"costs"`
	CostRatio           float64             `json:"cost_ratio"`
	CostRatioOperator   string              `json:"cost_ratio_operator"`
	Engine              string              `json:"engine"`
	Type                string              `json:"type"`
	Name                string              `json:"name"`
	StaticUtmSource     string              `json:"static_utm_source"`
	StaticUtmMedium     string              `json:"static_utm_medium"`
	StaticUtmCampaign   string              `json:"static_utm_campaign"`
	StaticUtmTerm       string              `json:"static_utm_term"`
	StaticUtmContent    string              `json:"static_utm_content"`
	StaticUtmReferrer   string              `json:"static_utm_referrer"`
	StaticUtmExpid      string              `json:"static_utm_expid"`
	DynamicCallTracking DynamicCallTracking `json:"dynamic_call_tracking"`
	CampaignConditions  []GroupCondition
	DateUpdate          time.Time
}
