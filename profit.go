package main

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type CryptoStats struct {
	NetworkHashRate float64
	AvgBlockTime    float64
	BlockReward     float64
}

var zcashStats CryptoStats
var ethereumStats CryptoStats
var bitcoinStats CryptoStats

func calculateExepectedProfit(hashRate, networkHashRate, avgBlockTime, period, reward float64) float64 {
	return hashRate / networkHashRate * period / avgBlockTime * reward
}

func getEthereumStats() {
	var currEthereumStats CryptoStats

	currEthereumStats.BlockReward = 5
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://etherscan.io/chart/blocktime?output=csv")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Ethereum block time: ", err)
		return
	}

	csv := strings.Split(string(resp.Body()), "\n")
	strTime := strings.Replace(strings.Replace(strings.Split(csv[len(csv)-2], ",")[2], "\"", "", -1), "\r", "", -1)
	currEthereumStats.AvgBlockTime, err = strconv.ParseFloat(strTime, 64)
	if err != nil {
		log.Println("Err parsing ethereum block time: ", err)
		return
	}

	req.SetRequestURI("https://etherscan.io/chart/hashrate?output=csv")

	resp = fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Ethereum hashrate: ", err)
		return
	}

	csv = strings.Split(string(resp.Body()), "\n")
	strHashrate := strings.Replace(strings.Replace(strings.Split(csv[len(csv)-2], ",")[2], "\"", "", -1), "\r", "", -1)
	currEthereumStats.NetworkHashRate, err = strconv.ParseFloat(strHashrate, 64)
	currEthereumStats.NetworkHashRate *= 1000000000
	if err != nil {
		log.Println("Err parsing ethereum block time: ", err)
		return
	}

	ethereumStats = currEthereumStats
}

func getZCashStats() {
	var currZCashStats CryptoStats

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://api.zcha.in/v2/mainnet/network")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting ZCash stats: ", err)
		return
	}

	body := resp.Body()
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Err unmarshaling zcash stats: ", err)
		return
	}

	currZCashStats.BlockReward = 10
	currZCashStats.AvgBlockTime = data["meanBlockTime"].(float64)
	currZCashStats.NetworkHashRate = data["hashrate"].(float64)

	zcashStats = currZCashStats
}

func getBitcoinStats() {
	var currBitcoinStats CryptoStats

	currBitcoinStats.BlockReward = 12.5
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://api.blockchain.info/charts/hash-rate?timespan=1year&format=json")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Bitcoin hashrate: ", err)
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Println("Err on unmarshaling bitcoin hashrate json")
		return
	}
	vals := data["values"].([]interface{})
	currBitcoinStats.NetworkHashRate = vals[len(vals)-1].(map[string]interface{})["y"].(float64) * 1000000000000
	req.SetRequestURI("https://api.blockchain.info/charts/difficulty?format=json")
	resp = fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Bitcoin difficulty: ", err)
		return
	}

	data = make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Println("Err on unmarshaling bitcoin hashrate json")
		return
	}
	vals = data["values"].([]interface{})
	difficulty := vals[len(vals)-1].(map[string]interface{})["y"].(float64)
	currBitcoinStats.AvgBlockTime = difficulty * math.Pow(2, 32) / currBitcoinStats.NetworkHashRate

	bitcoinStats = currBitcoinStats
}

func calculateZCashProfit(hashRate, period float64) float64 {
	return calculateExepectedProfit(hashRate, zcashStats.NetworkHashRate, zcashStats.AvgBlockTime,
		period, zcashStats.BlockReward)
}

func calculateEthereumProfit(hashRate, period float64) float64 {
	return calculateExepectedProfit(hashRate, ethereumStats.NetworkHashRate, ethereumStats.AvgBlockTime,
		period, ethereumStats.BlockReward)
}

func calculateBitcoinProfit(hashRate, period float64) float64 {
	return calculateExepectedProfit(hashRate, bitcoinStats.NetworkHashRate, bitcoinStats.AvgBlockTime,
		period, bitcoinStats.BlockReward)
}
