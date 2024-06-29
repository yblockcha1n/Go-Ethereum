package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	InfuraURL  string `json:"infura_url"`
	PrivateKey string `json:"private_key"`
	ToAddress  string `json:"to_address"`
}

func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := ethclient.Dial(config.InfuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve account nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	to := common.HexToAddress(config.ToAddress)
	value := big.NewInt(1000000000) // 1 gwei

	fmt.Print("Enter your message: ")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read message: %v", err)
	}
	message = strings.TrimSpace(message)

	// メッセージをHEXに変換
	data := []byte(message)
	hexData := hex.EncodeToString(data)
	fmt.Printf("Message in HEX: 0x%s\n", hexData)

	// HEXデータをスライスに変換
	txData := common.FromHex(hexData)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &to,
		Data: txData,
	})
	if err != nil {
		log.Fatalf("Failed to estimate gas limit: %v", err)
	}
	gasLimit += 30000　// テストとして余裕を持たせる

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, txData)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	receipt, err := waitForTransaction(ctx, client, signedTx.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	fmt.Printf("Transaction confirmed in block %d\n", receipt.BlockNumber)
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)
}

func waitForTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			fmt.Println("Waiting for transaction confirmation...")
		}
	}
}
