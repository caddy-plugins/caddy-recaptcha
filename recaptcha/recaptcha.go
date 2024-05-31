package recaptcha

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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

var errInvalid = errors.New("Failed to validate reCAPTCHA.")

func (h Recaptchas) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var append string
	isGet := r.Method == http.MethodGet
	for _, rule := range h.Rules {
		if r.URL.Path != rule.GetPath() {
			continue
		}
		if r.Method != rule.GetMethod() {
			if isGet {
				append += buildJSScript(rule.GetSiteKey(), rule.GetAction())
			}
			continue
		}
		if !rule.Validate(r) {
			return 400, errInvalid
		}
		//fmt.Println(`Success!!!!!!!!!!!!!`)
		return h.Next.ServeHTTP(w, r)
	}

	if isGet && len(append) > 0 {
		return filter.RewriteResponse(w, r, func(wrapper http.ResponseWriter) bool {
			return strings.SplitN(wrapper.Header().Get(`Content-Type`), `;`, 2)[0] == `text/html`
		}, func(wrapper http.ResponseWriter) (bool, []byte) {
			wrapperd := wrapper.(RecordedAndDecodeIfRequired)
			body := wrapperd.RecordedAndDecodeIfRequired()
			bodyRetrieved := true
			fmt.Println(string(body))
			body = fixedForm(body, append)
			return bodyRetrieved, body
		}, 1024*1024, h.Next)
	}
	return h.Next.ServeHTTP(w, r)
}

type RecordedAndDecodeIfRequired interface {
	RecordedAndDecodeIfRequired() []byte
}
