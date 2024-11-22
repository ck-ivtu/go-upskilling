package su4

import (
	"encoding/json"
	"net/http"
	"os"
)

type Coin struct {
	CoinId   string  `json:"coinId"`
	Coin     string  `json:"coin"`
	Transfer string  `json:"transfer"`
	Chains   []Chain `json:"chains"`
}

type Chain struct {
	Chain        string `json:"chain"`
	Withdrawable string `json:"withdrawable"`
	WithdrawFee  string `json:"withdrawFee"`
}

func ParseMarketData() {
	response, err := http.Get("https://api.bitget.com/api/spot/v1/public/currencies")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	outputFile, err := os.Create("coins_filtered.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	decoder := json.NewDecoder(response.Body)

	var apiResponse struct {
		Code string `json:"code"`
		Data []Coin `json:"data"`
	}

	if err := decoder.Decode(&apiResponse); err != nil {
		panic(err)
	}

	for _, coin := range apiResponse.Data {
		if len(coin.Chains) == 0 {
			continue
		}

		if coin.Transfer != "true" {
			continue
		}

		var filteredChains []Chain

		for _, chain := range coin.Chains {
			if chain.Withdrawable == "true" {
				filteredChains = append(filteredChains, chain)
			}
		}

		if len(filteredChains) == 0 {
			continue
		}

		coin.Chains = filteredChains

		coinJson, err := json.Marshal(coin)
		if err != nil {
			panic(err)
		}

		_, err = outputFile.Write(append(coinJson, '\n'))
		if err != nil {
			panic(err)
		}
	}
}
