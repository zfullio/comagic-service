package comagic

import (
	"fmt"
	"net/http"
)

type Client struct {
	host    Host
	Version string
	Token   string
	xSid    string
	tr      http.Client
}

func NewClient(host Host, version string, token string) *Client {
	return &Client{
		host:    host,
		Version: version,
		Token:   token,
		tr:      http.Client{},
	}
}

type Payload struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  Method      `json:"method"`
	Params  interface{} `json:"params"`
}

func (c *Client) NewPayload(method Method, params interface{}) *Payload {
	return &Payload{
		Jsonrpc: c.Version,
		Id:      generateID(method),
		Method:  method,
		Params:  params,
	}
}

type Response struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      string    `json:"id"`
	Error   RespError `json:"error"`
}
type RespError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Mnemonic string `json:"mnemonic"`
		Field    string `json:"field"`
		Value    string `json:"value"`
		Params   struct {
			Object string `json:"object"`
		} `json:"params"`
		ExtendedHelper string `json:"extended_helper"`
		Metadata       struct {
		} `json:"metadata"`
	} `json:"data"`
}

type Metadata struct {
	TotalItems int         `json:"total_items"`
	Version    interface{} `json:"version"`
	Limits     Limits      `json:"limits"`
}

type Limits struct {
	DayLimit        int `json:"day_limit"`
	DayRemaining    int `json:"day_remaining"`
	DayReset        int `json:"day_reset"`
	MinuteLimit     int `json:"minute_limit"`
	MinuteRemaining int `json:"minute_remaining"`
	MinuteReset     int `json:"minute_reset"`
}

func (e *RespError) Error() string {
	return fmt.Sprintf("%v:%s: API error", e.Code, e.Message)
}

type Method string

const (
	GetAccount               Method = "get.account"
	GetCampaigns             Method = "get.campaigns"
	GetVirtualNumbers        Method = "get.virtual_numbers"
	GetSites                 Method = "get.sites"
	GetCallsReport           Method = "get.calls_report"
	GetOfflineMessagesReport Method = "get.offline_messages_report"
)
