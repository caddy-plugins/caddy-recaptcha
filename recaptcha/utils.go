package recaptcha

import (
	"bytes"
	"slices"
	"strings"
)

func buildJSScript(sitekey string, action string) string {
	rpl := strings.NewReplacer(
		`{{jslib}}`, jslib,
		`{{sitekey}}`, sitekey,
		`{{action}}`, action,
	)
	return rpl.Replace(jsscript)
}

func fixedForm(content []byte, append string) []byte {
	index := bytes.LastIndex(content, []byte(`</body>`))
	if index <= 0 {
		return content
	}
	content = slices.Insert(content, index-1, []byte(append)...)
	return content
}
