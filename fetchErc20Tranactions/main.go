package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type EtherscanReponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		From              string `json:"from"`
		ContractAddress   string `json:"contractAddress"`
		To                string `json:"to"`
		Value             string `json:"value"`
		TokenName         string `json:"tokenName"`
		TokenSymbol       string `json:"tokenSymbol"`
		TokenDecimal      string `json:"tokenDecimal"`
		TransactionIndex  string `json:"transactionIndex"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		GasUsed           string `json:"gasUsed"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		Input             string `json:"input"`
		Confirmations     string `json:"confirmations"`
	} `json:"result"`
}

var (
	apiKey      = " <YOUR API KEY>"
	ourAddress  = "0x5DD596C901987A2b28C38A9C1DfBf86fFFc15d77"
	usdcContact = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
)

func main() {
	ether := "https://api.etherscan.io/api"

	var res *http.Response

	endpoint, err := url.Parse(ether)
	if err != nil {
		log.Fatal(err)
	}

	params := endpoint.Query()

	params.Set("page", "1")
	params.Set("sort", "asc")
	params.Set("offset", "100")
	params.Set("apikey", apiKey)
	params.Set("action", "tokentx")
	params.Set("module", "account")
	params.Set("address", ourAddress)
	params.Set("endblock", "14711463")
	params.Set("startblock", "14711461")
	params.Set("contractaddress", usdcContact)

	endpoint.RawQuery = params.Encode()

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		log.Println(err)
	}

	res, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var apiResponse EtherscanReponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(AsPrettyJson(apiResponse)))
}

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonB)
}
