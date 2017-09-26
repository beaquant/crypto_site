package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/muller95/tntsessions"
	"github.com/valyala/fasthttp"
)

type MainPage struct {
	Title           string
	CalculatorTitle string
	PerDay          string
	PerWeek         string
	PerMonth        string
	PerYear         string
	Calculate       string
	Language        string
	Simple          string
	EthereumYear    string
}

func mainPage(ctx *fasthttp.RequestCtx, sess *tntsessions.Session) {
	var mainPage MainPage
	bytes, err := ioutil.ReadFile("public/resources/" + sess.GetString("language") + "/main.json")
	if err != nil {
		log.Printf("Err reading public/resources/%v/main.json: %v\n", sess.GetString("language"), err)
		ctx.Response.SetStatusCode(int(InternalServerError))
		return
	}

	err = json.Unmarshal(bytes, &mainPage)
	if err != nil {
		log.Printf("Err unmarshaling public/resources/%v/main.json: %v\n",
			sess.GetString("language"), err)
		ctx.Response.SetStatusCode(int(InternalServerError))
		return
	}

	mainPage.Language = sess.GetString("language")

	template, err := template.ParseFiles("public/pages/main.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
	}
	ctx.SetContentType("text/html")
	err = template.Execute(ctx, mainPage)
	if err != nil {
		log.Println("Err on executing main page template: ", err)
		ctx.SetStatusCode(int(InternalServerError))
		return
	}

	ctx.SetStatusCode(int(Ok))
}
