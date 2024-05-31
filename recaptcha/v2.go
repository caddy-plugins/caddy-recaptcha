package recaptcha

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
)

type V2Rule struct {
	Secret  string
	Method  string
	Path    string
	SiteKey string
}

func (rule V2Rule) GetPath() string {
	return rule.Path
}

func (rule V2Rule) GetMethod() string {
	return rule.Method
}

func (rule V2Rule) GetAction() string {
	return ``
}

func (rule V2Rule) GetSiteKey() string {
	return rule.SiteKey
}

func (rule V2Rule) Validate(r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		return false
	}

	response := r.PostForm.Get("g-recaptcha-response")
	if len(response) == 0 {
		response = r.Header.Get("g-recaptcha-response")
	}
	if len(response) == 0 {
		return false
	}

	params := url.Values{}
	params.Set("secret", rule.Secret)
	params.Set("response", response)
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		return false
	}

	result := &V2Result{}
	json.NewDecoder(resp.Body).Decode(result)

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return false
	}

	if result.Hostname != host {
		return false
	}

	return result.Success
}

type V2Result struct {
	Success  bool   `json:"success"`
	Hostname string `json:"hostname"`
}
