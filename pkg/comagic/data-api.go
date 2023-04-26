package comagic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (c *Client) GetAccount() (data RespGetAccount, err error) {
	params := GetAccountParams{AccessToken: c.Token}
	payload := c.NewPayload(GetAccount, params)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return data, fmt.Errorf("ошибка формирования запроса: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return data, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.tr.Do(req)
	if err != nil {
		return data, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return data, err
	}

	if data.Error.Code != 0 {
		return data, &data.Error
	}

	return data, err

}

type GetAccountParams struct {
	AccessToken string `json:"access_token"`
}

type RespGetAccount struct {
	Response
	Result struct {
		Metadata Metadata `json:"metadata"`
		Data     []struct {
			AppId    int64  `json:"app_id"`
			Name     string `json:"name"`
			Timezone string `json:"timezone"`
		} `json:"data"`
	} `json:"result"`
}

func (c *Client) GetCampaigns(fields []string, filter Filter) (data RespCampaignsInfo, err error) {
	params := GetRequestParams{
		AccessToken: c.Token,
		Limit:       10000,
		Filter:      &filter,
		Fields:      fields,
	}
	payload := c.NewPayload(GetCampaigns, params)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return data, fmt.Errorf("ошибка формирования запроса: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return data, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.tr.Do(req)
	if err != nil {
		return data, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return data, err
	}

	if data.Error.Code != 0 {
		return data, &data.Error
	}

	return data, err

}

type Filter struct {
	Filters   []FilterParams `json:"filters,omitempty"`
	Condition string         `json:"condition,omitempty"`
}

type FilterParams struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

type SortParams struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type GetRequestParams struct {
	AccessToken string       `json:"access_token"`
	UserId      int64        `json:"user_id,omitempty"`
	Offset      int          `json:"offset,omitempty"`
	Limit       int          `json:"limit,omitempty"`
	Filter      *Filter      `json:"filter,omitempty"`
	Sort        []SortParams `json:"sort,omitempty"`
	Fields      []string     `json:"fields,omitempty"`
}

type RespCampaignsInfo struct {
	Response
	Result struct {
		Metadata Metadata       `json:"metadata"`
		Data     []CampaignInfo `json:"data"`
	} `json:"result"`
}

type CampaignInfo struct {
	ID                  int64               `json:"id"`
	Status              string              `json:"status"`
	CreationTime        string              `json:"creation_time"`
	Description         string              `json:"description"`
	SiteId              int64               `json:"site_id"`
	SiteDomainName      string              `json:"site_domain_name"`
	Costs               float64             `json:"costs"`
	CostRatio           float64             `json:"cost_ratio"`
	CostRatioOperator   string              `json:"cost_ratio_operator"`
	Engine              string              `json:"engine"`
	Type                string              `json:"type"`
	Name                string              `json:"name"`
	SiteBlocks          []SiteBlock         `json:"site_blocks"`
	StaticUtmSource     string              `json:"static_utm_source"`
	StaticUtmMedium     string              `json:"static_utm_medium"`
	StaticUtmCampaign   string              `json:"static_utm_campaign"`
	StaticUtmTerm       string              `json:"static_utm_term"`
	StaticUtmContent    string              `json:"static_utm_content"`
	StaticUtmReferrer   string              `json:"static_utm_referrer"`
	StaticUtmExpid      string              `json:"static_utm_expid"`
	DynamicCallTracking DynamicCallTracking `json:"dynamic_call_tracking"`
	CampaignConditions  struct {
		GroupConditions [][]GroupCondition `json:"group_conditions,omitempty"`
	} `json:"campaign_conditions"`
}

type GroupCondition struct {
	Type              string `json:"type"`
	Value             string `json:"value"`
	Operator          string `json:"operator"`
	CampaignParameter string `json:"campaign_parameter"`
}

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

type DynamicCallTracking struct {
	ReservationTime     string  `json:"reservation_time"`
	CountVirtualNumbers int64   `json:"count_virtual_numbers"`
	CountVisits         int64   `json:"count_visits"`
	CoverageVisitors    float64 `json:"coverage_visitors"`
}

func (c *Client) GetVirtualNumbers() (data RespVirtualNumbersInfo, err error) {
	params := GetRequestParams{
		AccessToken: c.Token,
	}
	payload := c.NewPayload(GetVirtualNumbers, params)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return data, fmt.Errorf("ошибка формирования запроса: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return data, err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := c.tr.Do(req)
	if err != nil {
		return data, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return data, err
	}

	if data.Error.Code != 0 {
		return data, &data.Error
	}

	return data, err

}

type RespVirtualNumbersInfo struct {
	Response
	Result struct {
		Metadata Metadata `json:"metadata"`
		Data     []struct {
			Id                     int64  `json:"id"`
			VirtualPhoneNumber     string `json:"virtual_phone_number"`
			RedirectionPhoneNumber string `json:"redirection_phone_number"`
			ActivationDate         string `json:"activation_date"`
			Status                 string `json:"status"`
			Category               string `json:"category"`
			Type                   string `json:"type"`
			Campaigns              []struct {
				CampaignId     int64  `json:"campaign_id"`
				SiteId         int64  `json:"site_id"`
				SiteDomainName string `json:"site_domain_name"`
				CampaignName   string `json:"campaign_name"`
				SiteBlocks     []struct {
					SiteBlockId   int64  `json:"site_block_id"`
					SiteBlockName string `json:"site_block_name"`
					IsTracking    bool   `json:"is_tracking"`
				} `json:"site_blocks"`
			} `json:"campaigns"`
			Scenarios []struct {
				ScenarioId   int64  `json:"scenario_id"`
				ScenarioName string `json:"scenario_name"`
			} `json:"scenarios"`
		} `json:"data"`
	} `json:"result"`
}

func (c *Client) GetSites() error {
	params := GetRequestParams{
		AccessToken: c.Token,
	}
	payload := c.NewPayload(GetSites, params)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка формирования запроса: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.tr.Do(req)
	if err != nil {
		return err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	data := RespSitesInfo{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return err
	}
	if data.Error.Code != 0 {
		return &data.Error
	}
	return nil

}

type RespSitesInfo struct {
	Response
	Result struct {
		Metadata Metadata `json:"metadata"`
		Data     []struct {
			Id                 int64  `json:"id"`
			DomainName         string `json:"domain_name"`
			Name               string `json:"name"`
			DefaultPhoneNumber string `json:"default_phone_number"`
			DefaultScenario    struct {
				ScenarioId   int64  `json:"scenario_id"`
				ScenarioName string `json:"scenario_name"`
			} `json:"default_scenario"`
			SiteKey                          string `json:"site_key"`
			IndustryId                       int64  `json:"industry_id"`
			IndustryName                     string `json:"industry_name"`
			TargetCallMinDuration            int64  `json:"target_call_min_duration"`
			TrackSubdomains                  bool   `json:"track_subdomains"`
			CookieLifetime                   int64  `json:"cookie_lifetime"`
			CampaignLifetime                 int64  `json:"campaign_lifetime"`
			SalesEnabled                     bool   `json:"sales_enabled"`
			SecondCommunicationPeriod        int64  `json:"second_communication_period"`
			ServicesEnabled                  bool   `json:"services_enabled"`
			ReplacementDynamicalBlockEnabled bool   `json:"replacement_dynamical_block_enabled"`
			WidgetLink                       struct {
				Enabled bool   `json:"enabled"`
				Text    string `json:"text"`
				Url     string `json:"url"`
			} `json:"widget_link"`
			ShowVisitorId struct {
				Enabled         bool   `json:"enabled"`
				ElementIdValue  string `json:"element_id_value"`
				Message         string `json:"message"`
				LengthVisitorId int64  `json:"length_visitor_id"`
			} `json:"show_visitor_id"`
			SiteBlocks []struct {
				SiteBlockId   int64  `json:"site_block_id"`
				SiteBlockName string `json:"site_block_name"`
			} `json:"site_blocks"`
			ConnectedIntegrations []interface{} `json:"connected_integrations"`
		} `json:"data"`
	} `json:"result"`
}

func (c *Client) GetCallsReport(dateFrom time.Time, dateTill time.Time, fields []string) (calls []CallInfo, err error) {
	if dateFrom.After(dateTill) {
		return calls, fmt.Errorf("дата окончания не может быть раньше даты начала")
	}

	receivedPositions := 0
	limit := 10000
	offset := 0
	for true {
		params := GetCallsRequestParams{
			GetRequestParams: GetRequestParams{
				AccessToken: c.Token,
				Limit:       limit,
				Fields:      fields,
				Offset:      offset,
			},
			DateFrom: dateFrom.Format(time.DateTime),
			DateTill: dateTill.Format(time.DateTime),
		}
		payload := c.NewPayload(GetCallsReport, params)
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return calls, fmt.Errorf("ошибка формирования запроса: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
		if err != nil {
			return calls, err
		}

		req.Header.Add("Content-Type", "application/json")
		resp, err := c.tr.Do(req)
		if err != nil {
			return calls, err
		}

		if resp.StatusCode != http.StatusOK {
			return calls, fmt.Errorf("ошибка получения ответа от API: %s", resp.Status)
		}

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return calls, err
		}

		var data = RespCallsReport{}
		data.Result.Data = make([]CallInfo, 0, limit)
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			return calls, err
		}

		if data.Error.Code != 0 {
			return calls, &data.Error
		}

		calls = append(calls, data.Result.Data...)
		receivedPositions = len(calls)
		if receivedPositions < data.Result.Metadata.TotalItems {
			ControlLimits(data.Result.Metadata.Limits)
			offset += limit
		} else {
			break
		}
	}
	return calls, err
}

func ControlLimits(l Limits) {
	if l.MinuteLimit <= 5 {
		log.Printf("Упреждение в минутный лимит. Пауза %v секунд\n", l.MinuteReset)
		time.Sleep(time.Duration(l.MinuteReset) * time.Second)
	}
}

func (c *Client) GetOfflineMessagesReport(dateFrom time.Time, dateTill time.Time, fields []string) (messages []OfflineMessageInfo, err error) {
	if dateFrom.After(dateTill) {
		return messages, fmt.Errorf("дата окончания не может быть раньше даты начала")
	}

	receivedPositions := 0
	limit := 10000
	offset := 0
	for true {
		params := GetCallsRequestParams{
			GetRequestParams: GetRequestParams{
				AccessToken: c.Token,
				Limit:       limit,
				Fields:      fields,
				Offset:      offset,
			},
			DateFrom: dateFrom.Format(time.DateTime),
			DateTill: dateTill.Format(time.DateTime),
		}
		payload := c.NewPayload(GetOfflineMessagesReport, params)
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return messages, fmt.Errorf("ошибка формирования запроса: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, c.buildLink(), bytes.NewBuffer(payloadJSON))
		if err != nil {
			return messages, err
		}

		req.Header.Add("Content-Type", "application/json")
		resp, err := c.tr.Do(req)
		if err != nil {
			return messages, err
		}

		if resp.StatusCode != http.StatusOK {
			return messages, fmt.Errorf("ошибка получения ответа от API: %s", resp.Status)
		}

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return messages, err
		}

		var data = RespOfflineMessagesReport{}
		data.Result.Data = make([]OfflineMessageInfo, 0, limit)
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			return messages, err
		}

		if data.Error.Code != 0 {
			return messages, &data.Error
		}

		messages = append(messages, data.Result.Data...)
		receivedPositions = len(messages)
		if receivedPositions < data.Result.Metadata.TotalItems {
			ControlLimits(data.Result.Metadata.Limits)
			offset += limit
		} else {
			break
		}
	}
	return messages, err
}

type GetCallsRequestParams struct {
	GetRequestParams
	DateFrom string `json:"date_from"`
	DateTill string `json:"date_till"`
}

type RespCallsReport struct {
	Response
	Result struct {
		Data     []CallInfo `json:"data"`
		Metadata Metadata   `json:"metadata"`
	} `json:"result"`
}

type RespOfflineMessagesReport struct {
	Response
	Result struct {
		Data     []OfflineMessageInfo `json:"data"`
		Metadata Metadata             `json:"metadata"`
	} `json:"result"`
}

type CallInfo struct {
	Id                            int64            `json:"id"`
	StartTime                     string           `json:"start_time"`
	FinishTime                    string           `json:"finish_time"`
	VirtualPhoneNumber            string           `json:"virtual_phone_number"`
	IsTransfer                    bool             `json:"is_transfer"`
	FinishReason                  string           `json:"finish_reason"`
	Direction                     string           `json:"direction"`
	Source                        string           `json:"source"`
	CommunicationNumber           int64            `json:"communication_number"`
	CommunicationPageUrl          string           `json:"communication_page_url"`
	CommunicationId               int64            `json:"communication_id"`
	CommunicationType             string           `json:"communication_type"`
	IsLost                        bool             `json:"is_lost"`
	CpnRegionId                   int64            `json:"cpn_region_id"`
	CpnRegionName                 string           `json:"cpn_region_name"`
	WaitDuration                  int64            `json:"wait_duration"`
	TotalWaitDuration             int64            `json:"total_wait_duration"`
	LostCallProcessingDuration    int64            `json:"lost_call_processing_duration"`
	TalkDuration                  int64            `json:"talk_duration"`
	CleanTalkDuration             int64            `json:"clean_talk_duration"`
	TotalDuration                 int64            `json:"total_duration"`
	PostprocessDuration           int64            `json:"postprocess_duration"`
	UaClientId                    string           `json:"ua_client_id"`
	YmClientId                    string           `json:"ym_client_id"`
	SaleDate                      string           `json:"sale_date"`
	SaleCost                      float64          `json:"sale_cost"`
	SearchQuery                   string           `json:"search_query"`
	SearchEngine                  string           `json:"search_engine"`
	ReferrerDomain                string           `json:"referrer_domain"`
	Referrer                      string           `json:"referrer"`
	EntrancePage                  string           `json:"entrance_page"`
	Gclid                         string           `json:"gclid"`
	Yclid                         string           `json:"yclid"`
	Ymclid                        string           `json:"ymclid"`
	EfId                          string           `json:"ef_id"`
	Channel                       string           `json:"channel"`
	SiteId                        int64            `json:"site_id"`
	SiteDomainName                string           `json:"site_domain_name"`
	CampaignId                    int64            `json:"campaign_id"`
	CampaignName                  string           `json:"campaign_name"`
	AutoCallCampaignName          string           `json:"auto_call_campaign_name"`
	VisitOtherCampaign            bool             `json:"visit_other_campaign"`
	VisitorId                     int64            `json:"visitor_id"`
	PersonId                      int64            `json:"person_id"`
	VisitorType                   string           `json:"visitor_type"`
	VisitorSessionId              int64            `json:"visitor_session_id"`
	VisitsCount                   int64            `json:"visits_count"`
	VisitorFirstCampaignId        int64            `json:"visitor_first_campaign_id"`
	VisitorFirstCampaignName      string           `json:"visitor_first_campaign_name"`
	VisitorCity                   string           `json:"visitor_city"`
	VisitorRegion                 string           `json:"visitor_region"`
	VisitorCountry                string           `json:"visitor_country"`
	VisitorDevice                 string           `json:"visitor_device"`
	LastAnsweredEmployeeId        int64            `json:"last_answered_employee_id"`
	LastAnsweredEmployeeFullName  string           `json:"last_answered_employee_full_name"`
	LastAnsweredEmployeeRating    int64            `json:"last_answered_employee_rating"`
	FirstAnsweredEmployeeId       int64            `json:"first_answered_employee_id"`
	FirstAnsweredEmployeeFullName string           `json:"first_answered_employee_full_name"`
	ScenarioId                    int64            `json:"scenario_id"`
	ScenarioName                  string           `json:"scenario_name"`
	CallApiExternalId             string           `json:"call_api_external_id"`
	CallApiRequestId              int64            `json:"call_api_request_id"`
	ContactPhoneNumber            string           `json:"contact_phone_number"`
	ContactFullName               string           `json:"contact_full_name"`
	ContactId                     int64            `json:"contact_id"`
	UtmSource                     string           `json:"utm_source"`
	UtmMedium                     string           `json:"utm_medium"`
	UtmTerm                       string           `json:"utm_term"`
	UtmContent                    string           `json:"utm_content"`
	UtmCampaign                   string           `json:"utm_campaign"`
	OpenstatAd                    string           `json:"openstat_ad"`
	OpenstatCampaign              string           `json:"openstat_campaign"`
	OpenstatService               string           `json:"openstat_service"`
	OpenstatSource                string           `json:"openstat_source"`
	EqUtmSource                   string           `json:"eq_utm_source"`
	EqUtmMedium                   string           `json:"eq_utm_medium"`
	EqUtmTerm                     string           `json:"eq_utm_term"`
	EqUtmContent                  string           `json:"eq_utm_content"`
	EqUtmCampaign                 string           `json:"eq_utm_campaign"`
	EqUtmReferrer                 string           `json:"eq_utm_referrer"`
	EqUtmExpid                    string           `json:"eq_utm_expid"`
	Attributes                    []string         `json:"attributes"`
	CallRecords                   []string         `json:"call_records"`
	VoiceMailRecords              []string         `json:"voice_mail_records"`
	Tags                          []Tag            `json:"tags"`
	VisitorCustomProperties       []CustomProperty `json:"visitor_custom_properties"`
	Segments                      []Segment        `json:"segments"`
	Employees                     []Employee       `json:"employees"`
	ScenarioOperations            []struct {
		Name string `json:"name"`
		Id   int64  `json:"id"`
	} `json:"scenario_operations"`
}
type Tag struct {
	TagId               int64  `json:"tag_id"`
	TagName             string `json:"tag_name"`
	TagType             string `json:"tag_type"`
	TagUserId           int64  `json:"tag_user_id"`
	TagUserLogin        string `json:"tag_user_login"`
	TagChangeTime       string `json:"tag_change_time"`
	TagEmployeeId       int64  `json:"tag_employee_id"`
	TagEmployeeFullName string `json:"tag_employee_full_name"`
}
type CustomProperty struct {
	PropertyName  string `json:"property_name"`
	PropertyValue string `json:"property_value"`
}
type Employee struct {
	EmployeeId       int64  `json:"employee_id"`
	EmployeeFullName string `json:"employee_full_name"`
	IsAnswered       bool   `json:"is_answered"`
}
type Segment struct {
	SegmentId   int64  `json:"segment_id"`
	SegmentName string `json:"segment_name"`
}

type OfflineMessageInfo struct {
	Id                       int64            `json:"id"`
	DateTime                 string           `json:"date_time"`
	Text                     string           `json:"text"`
	CommunicationNumber      int64            `json:"communication_number"`
	CommunicationPageUrl     string           `json:"communication_page_url"`
	CommunicationType        string           `json:"communication_type"`
	CommunicationId          int64            `json:"communication_id"`
	UaClientId               string           `json:"ua_client_id"`
	YmClientId               string           `json:"ym_client_id"`
	SaleDate                 string           `json:"sale_date"`
	SaleCost                 float64          `json:"sale_cost"`
	Status                   string           `json:"status"`
	ProcessTime              string           `json:"process_time"`
	FormType                 string           `json:"form_type"`
	FormName                 string           `json:"form_name"`
	SearchQuery              string           `json:"search_query"`
	SearchEngine             string           `json:"search_engine"`
	ReferrerDomain           string           `json:"referrer_domain"`
	Referrer                 string           `json:"referrer"`
	EntrancePage             string           `json:"entrance_page"`
	Gclid                    string           `json:"gclid"`
	Yclid                    string           `json:"yclid"`
	Ymclid                   string           `json:"ymclid"`
	EfId                     string           `json:"ef_id"`
	Channel                  string           `json:"channel"`
	EmployeeId               int64            `json:"employee_id"`
	EmployeeFullName         string           `json:"employee_full_name"`
	EmployeeAnswerMessage    string           `json:"employee_answer_message"`
	EmployeeComment          string           `json:"employee_comment"`
	Tags                     []Tag            `json:"tags"`
	SiteId                   int64            `json:"site_id"`
	SiteDomainName           string           `json:"site_domain_name"`
	GroupId                  int64            `json:"group_id"`
	GroupName                string           `json:"group_name"`
	CampaignId               int64            `json:"campaign_id"`
	CampaignName             string           `json:"campaign_name"`
	VisitOtherCampaign       bool             `json:"visit_other_campaign"`
	VisitorId                int64            `json:"visitor_id"`
	VisitorName              string           `json:"visitor_name"`
	VisitorPhoneNumber       string           `json:"visitor_phone_number"`
	VisitorEmail             string           `json:"visitor_email"`
	PersonId                 int64            `json:"person_id"`
	VisitorType              string           `json:"visitor_type"`
	VisitorSessionId         int64            `json:"visitor_session_id"`
	VisitsCount              int64            `json:"visits_count"`
	VisitorFirstCampaignId   int64            `json:"visitor_first_campaign_id"`
	VisitorFirstCampaignName string           `json:"visitor_first_campaign_name"`
	VisitorCity              string           `json:"visitor_city"`
	VisitorRegion            string           `json:"visitor_region"`
	VisitorCountry           string           `json:"visitor_country"`
	VisitorDevice            string           `json:"visitor_device"`
	VisitorCustomProperties  []CustomProperty `json:"visitor_custom_properties"`
	Segments                 []Segment        `json:"segments"`
	UtmSource                string           `json:"utm_source"`
	UtmMedium                string           `json:"utm_medium"`
	UtmTerm                  string           `json:"utm_term"`
	UtmContent               string           `json:"utm_content"`
	UtmCampaign              string           `json:"utm_campaign"`
	OpenstatAd               string           `json:"openstat_ad"`
	OpenstatCampaign         string           `json:"openstat_campaign"`
	OpenstatService          string           `json:"openstat_service"`
	OpenstatSource           string           `json:"openstat_source"`
	EqUtmSource              string           `json:"eq_utm_source"`
	EqUtmMedium              string           `json:"eq_utm_medium"`
	EqUtmTerm                string           `json:"eq_utm_term"`
	EqUtmContent             string           `json:"eq_utm_content"`
	EqUtmCampaign            string           `json:"eq_utm_campaign"`
	EqUtmReferrer            string           `json:"eq_utm_referrer"`
	EqUtmExpid               string           `json:"eq_utm_expid"`
	Attributes               []string         `json:"attributes"`
	SourceId                 int64            `json:"source_id"`
	SourceName               string           `json:"source_name"`
	SourceNew                string           `json:"source_new"`
	ChannelNew               string           `json:"channel_new"`
	ChannelCode              string           `json:"channel_code"`
}
