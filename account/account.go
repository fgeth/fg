package account

import (
    "crypto/ecdsa"
    "fmt"
    "log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/block"

)

type Account struct {
	Address				string					//Public key as string
	Balance				big.Int					//Value of the account at the completion of the latest Block if Block 3 was just completed this is the value after all transactions have been completed in Block 3
	TxNumber			uint64					//Number of transactions sent from this account
	BlockNumber			uint64					//Block Number of current Balance
	Data				map[uint]BlockData	    //Map index is the year the data is valid for
		
}

type BlockData struct {
	Balance				map[uint64]big.Int		//Map index is Block Number and the associated value of the account at that Block Number. Balance[3] would show the value of the account at the end of Block 3 and would include any of block 3 transactions
	TxTo				[]common.Hash			//Array of transaction hashes that where sent to this account
	ValRew				[]ValidationReward		//Array of Rewards earned for validation
	TxFrom				[]common.Hash			//Array of transaction hashes that where sent from this account
	StartingBalance      big.Int				//Starting Balance for this year
}

type ValidationReward struct{
	BlockNumber		uint64
	Amount			big.Int
	PubKey			*ecdsa.PublicKey			
}

func NewAccount(password string, blockNumber uint64) Account{
 prvKey, err := createKey()
 pubKey := &prvKey.PublicKey

 privateKey, publicKey :=Encode(prvKey, pubKey)
 keyjson, err := Encrypt([]byte(password), []byte(privateKey))
	if err != nil {
		return err
	}
	balance := new big.Int(0)
	
	return account := new Account{publicKey, privateKey, pubKey, balance, balance, blockNumber} 

}

func SaveAccountToDisk(account Account){
	fileName := account.Address
   file, _ := json.MarshalIndent(account, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}

func LoadAccountFromDisk(address string) Account{

file, _ := ioutil.ReadFile(address)
 
	account := Account{}
 
	_ = json.Unmarshal([]byte(file), &account)

	return account
}

