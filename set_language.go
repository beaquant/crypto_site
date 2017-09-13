package main

import (
	"log"

	"github.com/muller95/tntsessions"
	"github.com/valyala/fasthttp"
)

func setLanguage(ctx *fasthttp.RequestCtx, sess *tntsessions.Session) {
	language := string(ctx.PostArgs().Peek("language"))
	if language != "ru" {
		language = "en"
	}

	sess.Set("language", language)
	err := sessDB.Put(sess)
	if err != nil {
		log.Printf("Err on setting language: %v\n", err)
		ctx.Response.SetStatusCode(int(InternalServerError))
	}

	ctx.Response.SetStatusCode(int(Ok))
}
