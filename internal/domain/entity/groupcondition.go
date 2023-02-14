package entity

import "time"

type GroupCondition struct {
	ID                int
	GroupID           int
	Type              string `json:"type"`
	Value             string `json:"value"`
	Operator          string `json:"operator"`
	CampaignParameter string `json:"campaign_parameter"`
	DateUpdate        time.Time
}
