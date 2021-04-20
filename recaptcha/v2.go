package recaptcha

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"path"
)

type V2Rule struct {
	Secret string
	Method string
	Path   string
}

func (rule V2Rule) Validate(r *http.Request) bool {
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
