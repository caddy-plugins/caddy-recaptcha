package main

import (
	"github.com/admpub/caddy/caddy/caddymain"
	_ "github.com/caddy-plugins/caddy-recaptcha"
)

func main() {
	caddymain.Run()
}
