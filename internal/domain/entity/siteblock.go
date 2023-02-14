package entity

type SiteBlock struct {
	SiteBlockId                int64  `json:"site_block_id"`
	SiteBlockName              string `json:"site_block_name"`
	PhoneNumberType            string `json:"phone_number_type"`
	PhoneNumberId              int64  `json:"phone_number_id"`
	PhoneNumber                string `json:"phone_number"`
	RedirectionPhoneNumberId   int64  `json:"redirection_phone_number_id"`
	RedirectionPhoneNumber     string `json:"redirection_phone_number"`
	DynamicCallTrackingEnabled bool   `json:"dynamic_call_tracking_enabled"`
}
