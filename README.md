# Go-Ethereum

https://github.com/yblockcha1n/Go-Ethereum/assets/144770048/aa763a8a-10c6-4c9f-8a0f-5638b0e85b54

## Prerequisites

- Go 1.16 or later
- Git

## Setup

**1. Clone the repository**
```
git clone https://github.com/yblockcha1n/Go-Ethereum.git
cd Go-Ethereum
```

**2. Install dependencies**
```
go mod init Go-Ethereum
go get github.com/ethereum/go-ethereum
```

**3. Create a `config.json` file in the project root with the following structure**

```json
{
  "infura_url": "https://arbitrum-sepolia.infura.io/v3/{YOUR_INFURA_ID}",
  "private_key": "YOUR_PRIVATE_KEY",
  "to_address": "SEND_TO_ADDRESS"
}
```

Replace the placeholders with your actual Infura URL, private key, and recipient address

## Running the Application

**1. Build the application**
```
go build -o Go-Ethereum
```

**2. Run the application**
```
./Go-Ethereum
```

**3. When prompted, enter the message you want to include in the transaction.
4. The application will convert your message to HEX, create a transaction including this data, and send it to the Sepolia Testnet.
5. Wait for the transaction to be confirmed. The application will display the transaction hash and block number once confirmed.**


## Important Notes

・Ensure you have sufficient ETH in your Sepolia Testnet account to cover gas fees.


・Never share your private key or commit it to version control.


・This application is for educational purposes and should not be used with real funds without proper security audits.

