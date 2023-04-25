package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/big"
	"strings"

	"example.xyz/m/v2/contracts/dai"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const (
	// nodeRPCUrl         = "https://cloudflare-eth.com/"
	GAS_LIMIT          = uint64(1_200_000)
	blockExplorBaseUrl = "https://rinkeby.etherscan.io"
	nodeRPCUrl         = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	privateKeyBytes    = "0xPRIVATEEEEEEEE Key"
)

var (
	wethContractAddress = common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab")
)

func main() {
	ethClient, err := ethclient.Dial(nodeRPCUrl)
	if err != nil {
		panic(err)
	}

	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println("Current Chain ID is 🔗 ", chainId)

	privateKey, err := crypto.HexToECDSA(privateKeyBytes[2:])
	if err != nil {
		log.Fatal("failed to parse PRIVATE_KEY", err.Error())
	}

	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := ethClient.PendingNonceAt(context.Background(), publicAddress)
	if err != nil {
		log.Fatal(err)
	}

	daiABI, err := abi.JSON(strings.NewReader(dai.DaiABI))
	if err != nil {
		panic(err)
	}

	val := math.Pow(10, 18)
	amount := int64(0.0001 * val)
	wethToSend := big.NewInt(amount)

	destinationAddr := common.HexToAddress("0xa76A2f0871CCf57521334B589fb9c9B6346D8856")

	rawTransData, err := daiABI.Pack("transfer", destinationAddr, wethToSend)
	if err != nil {
		panic(err)
	}

	log.Println("💀 Transaction Data: ", AsPrettyJson(rawTransData))

	gasFeeCap, err := ethClient.SuggestGasPrice(context.Background()) // new(big.Float).Mul(big.NewFloat(gasPrice.BlockPrices[0].BaseFeePerGas), new(big.Float).SetFloat64(math.Pow(10.0, 9.0))).Int(nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to get gas fee cap"))
	}

	gasTipCap, err := ethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed to get gas fee cap"))
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Gas:       GAS_LIMIT,
		Value:     common.Big0,
		Data:      rawTransData,
		ChainID:   chainId,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		To:        &wethContractAddress,
	})

	log.Println("🧾 Raw transaction : ", AsPrettyJson(tx))

	// We need a custom signer that can change the amout of gas it pays
	signer := types.NewLondonSigner(chainId)

	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		panic(errors.Wrap(err, "failed to sign transaction data during build"))
	}

	log.Println("🔐 Signed Transaction", AsPrettyJson(signedTx))

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	log.Printf("🌈 Commtted Transaction %s/tx/%s\n", blockExplorBaseUrl, signedTx.Hash())
}

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonB)
}
