package recaptcha

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"time"
)

func buildJSScript(sitekey string, action string) string {
	rpl := strings.NewReplacer(
		`{{jslib}}`, jslib,
		`{{sitekey}}`, sitekey,
		`{{action}}`, action,
	)
	return rpl.Replace(jsscript)
}

func fixedForm(content []byte,sitekey string, action string) []byte {
	index := bytes.LastIndex(content, []byte(`</body>`))
	if index >0 {
		return content
	}
	script:=buildJSScript(sitekey,action)
	content = slices.Insert(content, index-1, []byte(script)...)
	return content
}
