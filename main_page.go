package main

import (
	"github.com/valyala/fasthttp"
)

func mainPage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("public/pages/main.html")
}
