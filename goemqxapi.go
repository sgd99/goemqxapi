package goemqxapi

import (
	"encoding/base64"
)

type Goemq struct {
	BaseURL   string `json:"base_url"`
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func NewGoemq(baseURL, appId, appSecret string) *Goemq {
	return &Goemq{BaseURL: baseURL, AppId: appId, AppSecret: appSecret}
}

func (g *Goemq) getBasicAuthHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(g.AppId+":"+g.AppSecret))
}
