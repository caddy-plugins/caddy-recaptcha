package recaptcha

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"path"
)

type V3Rule struct {
	Secret    string
	Action    string
	Threshold float64
	Method    string
	Path      string
}

func (rule V3Rule) Validate(r *http.Request) bool {
	if r.Method != rule.Method {
		return true
	}

	if path.Clean(r.URL.Path) != rule.Path {
		return true
	}

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
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", params)
	if err != nil {
		return false
	}

	result := &V3Result{}
	json.NewDecoder(resp.Body).Decode(result)

	if !result.Success {
		return false
	}

	if result.Action != rule.Action {
		return false
	}

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return false
	}

	if result.Hostname != host {
		return false
	}

	if result.Score < rule.Threshold {
		return false
	}

	return true
}

type V3Result struct {
	Success  bool    `json:"success"`
	Action   string  `json:"action"`
	Hostname string  `json:"hostname"`
	Score    float64 `json:"score"`
}
