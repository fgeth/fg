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

)

type Account struct {
	Address				string
	Key					PrvKey
	Balance				big.Int					// Value of the account at the completion of the latest Block if Block 3 was just completed this is the value after all transactions have been completed in Block 3
	Pending				big.Int					//Amount after any known Transactions
	BlockNumber			uint64					//Block Number of the last Known completed Block
	Data				BlockData
	Confirmations		[]SignedTx				//Balance Confirmations sigend by each validator Node
}



func newAccount(password string) Account{


}

