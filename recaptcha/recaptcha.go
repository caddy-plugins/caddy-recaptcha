package recaptcha

import (
	"errors"
	"net/http"

	"github.com/admpub/caddy/caddyhttp/httpserver"
	filter "github.com/caddy-plugins/caddy-filter"
)

type Recaptchas struct {
	Next  httpserver.Handler
	Rules []Rule
}

type Rule interface {
	GetPath() string
	GetMethod() string
	GetAction() string
	GetSiteKey() string
	Validate(r *http.Request) bool
}

func (h Recaptchas) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var append string
	isGetHTML := r.Method == http.MethodGet && w.Header().Get(`Content-Type`) == `text/html`
	for _, rule := range h.Rules {
		if r.URL.Path != rule.GetPath() {
			continue
		}
		if r.Method != rule.GetMethod() {
			if isGetHTML {
				append += buildJSScript(rule.GetSiteKey(), rule.GetAction())
			}
			continue
		}
		if !rule.Validate(r) {
			return 400, errors.New("Failed to validate reCAPTCHA.")
		}
		return h.Next.ServeHTTP(w, r)
	}

	if isGetHTML && len(append) > 0 {
		return filter.RewriteResponse(w, r, func(wrapper http.ResponseWriter) bool {
			return true
		}, func(wrapper http.ResponseWriter) (bool, []byte) {
			wrapperd := wrapper.(RecordedAndDecodeIfRequired)
			body := wrapperd.RecordedAndDecodeIfRequired()
			bodyRetrieved := true
			body = fixedForm(body, append)
			return bodyRetrieved, body
		}, 1024*1024, h.Next)
	}
	return h.Next.ServeHTTP(w, r)
}

type RecordedAndDecodeIfRequired interface {
	RecordedAndDecodeIfRequired() []byte
}
