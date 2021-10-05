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
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/block"

)

type Account struct {
	Address				string
	Key					PrvKey
	Balance				big.Int					// Value of the account at the completion of the latest Block if Block 3 was just completed this is the value after all transactions have been completed in Block 3
	Pending				big.Int					//Amount after any known Transactions
	TxNumber			uint64					//Number of transactions sent from this account
	BlockNumber			uint64					//Block Number of the last Known completed Block
	Data				map[uint]BlockData	    //Map index is the year the data is valid for
	TxTo				[]uint64				//Array of transaction hashes that where sent to this account
	TxFrom				[]uint64				//Array of transaction hashes that where sent from this account
	ValRew				[]ValidationReward		//Array of Rewards earned for validation	
}

type BlockData struct {
	Year				uint					   //The year this data is valid for	
	Balance				map[uint64]big.Int		  //Map index is Block Number and the associated value of the account at that Block Number. History[3] would show the value of the account at the end of Block 2 and to include any of block 3 transactions
	Confirmations		map[uint64][]SignedTx	  //Map index is Block Number and the associated cofirmations of the account Balance at this height	
	EBLY				map[uint]big.Int		  //The Ending balance of past years with the map index equal to the year
}

type ValidationReward struct{
	BlockNumber		uint64
	Amount			big.Int
	PubKey			*ecdsa.PublicKey			
}

func NewAccount(password string) Account{


}

func SaveAccountToDisk(account Account){


}

func LoadAccountFromDisk(address string){


}

