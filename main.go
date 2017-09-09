package main

import (
	"log"

	"encoding/json"

	"time"

	"io/ioutil"

	"github.com/valyala/fasthttp"
)

var profitFunc = map[string]func(hashRate, period float64) float64{
	"/bitcoin_profit":  calculateBitcoinProfit,
	"/ethereum_profit": calculateEthereumProfit,
	"/zcash_profit":    calculateZCashProfit,
}

var gpuHashrates map[string]map[string]float64

func updateProfitRoutine() {
	time.Sleep(5 * time.Minute)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	switch path {
	case "/":
		mainPage(ctx)
	case "/bitcoin_profit", "/ethereum_profit", "/zcash_profit":
		hashrate := ctx.QueryArgs().GetUfloatOrZero("hashrate")
		perDay := profitFunc[path](hashrate, 24*60*60)
		perWeek := profitFunc[path](hashrate, 7*24*60*60)
		perMonth := profitFunc[path](hashrate, 30*24*60*60)
		perYear := profitFunc[path](hashrate, 365*24*60*60)

		profitMap := map[string]interface{}{
			"per_day":   perDay,
			"per_week":  perWeek,
			"per_month": perMonth,
			"per_year":  perYear,
		}

		profitJSON, _ := json.Marshal(profitMap)
		ctx.Write(profitJSON)
	default:
		fasthttp.FSHandler("public/", 0)(ctx)
	}
}

func main() {
	getEthereumStats()
	getZCashStats()
	getBitcoinStats()
	go updateProfitRoutine()

	bytes, _ := ioutil.ReadFile("data/hashrates.json")
	json.Unmarshal(bytes, &gpuHashrates)

	err := fasthttp.ListenAndServe("0.0.0.0:8080", requestHandler)
	// err = fasthttp.ListenAndServeTLS(":"+serverPort, certificatePath, keyPath,
	// requestHandler)
	if err != nil {
		log.Fatal("Err on startup server: ", err)
	}
}
