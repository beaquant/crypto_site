package main

import (
	"fmt"
	"log"
	"strings"

	"encoding/json"

	"time"

	"io/ioutil"

	"github.com/muller95/tntsessions"
	"github.com/valyala/fasthttp"
)

type RestCode uint32

const (
	Ok                  RestCode = 200
	NotFound            RestCode = 404
	SessionExpired      RestCode = 471
	InternalServerError RestCode = 500
)

var ethereumReward = 5.0

var profitFunc = map[string]func(hashRate, period float64) float64{
	"/bitcoin_profit":  calculateBitcoinProfit,
	"/ethereum_profit": calculateEthereumProfit,
	"/zcash_profit":    calculateZCashProfit,
}

var gpuHashrates map[string]map[string]float64
var sessDB *tntsessions.SessionsBase
var sessions map[string]*tntsessions.Session
var commitTime = 5
var updateProfit = 60
var ethereumCoefficients []float64
var ethereumPrices []float64

func updateProfitRoutine() {
	time.Sleep(time.Duration(updateProfit) * time.Minute)
	getEthereumStats()
	getZCashStats()
	getBitcoinStats()
}

func commitSessionsRoutine() {
	time.Sleep(time.Duration(commitTime) * time.Minute)
	for _, sess := range sessions {
		if sess.EndTime > time.Now().Unix() {
			delete(sessions, sess.ID)
			err := sessDB.Delete(sess.ID)
			if err != nil {
				log.Printf("Err deleting session %v: %v\n", sess, err)
			}
		}
		err := sessDB.Put(sess)
		if err != nil {
			log.Printf("Err putting session %v: %v\n", sess, err)
		}
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	var err error

	language := "en"
	acceptLanguages := strings.Split(string(ctx.Request.Header.Peek("Accept-Language")), ";")
	if len(acceptLanguages) > 0 {
		acceptLanguage := strings.Split(acceptLanguages[0], ",")
		if len(acceptLanguage) > 1 {
			if acceptLanguage[1] == "ru" {
				language = acceptLanguage[1]
			}
		}
	}

	sessID := string(ctx.Request.Header.Cookie("session_id"))
	sess, ok := sessions[sessID]
	if !ok {
		sess, err = sessDB.Get(sessID)
	}

	if err != nil && err != tntsessions.ErrNotFound {
		log.Printf("Err on getting session %v: %v\n", sessID, err)
		ctx.Response.SetStatusCode(int(InternalServerError))
	} else if err == tntsessions.ErrNotFound {
		sess, err = sessDB.Create(3 * 24 * 60 * 60)
		if err != nil {
			log.Printf("Err on creating session: %v\n", err)
			ctx.Response.SetStatusCode(int(InternalServerError))
			return
		}

		sessions[sess.ID] = sess
		sess.Set("language", language)
		sessions[sess.ID] = sess
		c := fasthttp.Cookie{}
		c.SetKey("session_id")
		c.SetValue(sess.ID)
		ctx.Response.Header.SetCookie(&c)
	}

	path := string(ctx.Path())
	switch path {
	case "/":
		mainPage(ctx, sess)
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
	case "/ethereum_prediction":
		// return hashRate / networkHashRate * period / avgBlockTime * reward

		hashrate := ctx.QueryArgs().GetUfloatOrZero("hashrate")
		powerConsumption := ctx.QueryArgs().GetUfloatOrZero("power_consumption")
		powerCost := ctx.QueryArgs().GetUfloatOrZero("power_cost")
		initialInvestment := ctx.QueryArgs().GetUfloatOrZero("initial_investment")
		log.Println(initialInvestment)
		profitEthereum := make([]float64, len(ethereumCoefficients))
		profitCurrency := make([]float64, len(ethereumPrices))
		for i := 0; i < len(ethereumCoefficients); i++ {
			profitEthereum[i] = 24 * 60 * 60 * hashrate * ethereumReward /
				ethereumCoefficients[i] / ethereumNetworkMultiplier
			// log.Printf("24*60*60*%f*%f*%f/%f=%f\n", hashrate, ethereumReward, ethereumCoefficients[i],
			// ethereumNetworkMultiplier, profitEthereum[i])

			profitCurrency[i] = profitEthereum[i]*ethereumPrices[i] -
				powerConsumption*powerCost*24.0/1000.0
			// log.Printf("price=%f", ethereumPrices[i])
			if i == 0 {
				profitCurrency[i] -= initialInvestment
			} else {
				profitEthereum[i] += profitEthereum[i-1]
				profitCurrency[i] += profitCurrency[i-1]
			}

		}

		profitJSON, _ := json.Marshal(map[string]interface{}{
			"ethereum": profitEthereum,
			"currency": profitCurrency,
		})
		ctx.Response.SetStatusCode(int(Ok))
		ctx.Write(profitJSON)
	case "/set_language":
		setLanguage(ctx, sess)
	default:
		fasthttp.FSHandler("public/", 0)(ctx)
	}
}

func main() {
	var err error

	predictEthereumParams()
	getEthereumStats()
	getZCashStats()
	getBitcoinStats()
	go updateProfitRoutine()

	bytes, _ := ioutil.ReadFile("data/hashrates.json")
	json.Unmarshal(bytes, &gpuHashrates)

	sessDB, err = tntsessions.ConnectToTarantool("127.0.0.1:3309", "guest", "", "sessions")
	if err != nil {
		log.Fatalf("Err on connecting to sessions db: %v\n", err)
	}

	sessions = make(map[string]*tntsessions.Session)

	fmt.Println("Server is ready")
	err = fasthttp.ListenAndServe("0.0.0.0:8080", requestHandler)

	// err = fasthttp.ListenAndServeTLS(":"+serverPort, certificatePath, keyPath,
	// requestHandler)
	if err != nil {
		log.Fatal("Err on startup server: ", err)
	}
}
