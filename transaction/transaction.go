package transaction

import (
    "crypto/sha1"
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"bufio"
	"hash"
	"math/big"
	"sync"
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	 "os"
	 "github.com/fgeth/fg/common"

)


type Transaction struct {
	TxId				string						//Year::Month::Day::Hour::Min::NodeID::TransactionNumber::Random4Digits
	BlockNumber			unit64						//Block Number transaction was confirmed
	From				address
	To					address
	Value				big.Int
	TxNumber			uint64						//The senders total number of sent transactions to this point including this transaction
	FVB					big.Int						//From Account Balance Value Before Transaction
	TVB					big.Int						//To Account Balance Value Before Transaction
	FVA					big.Int						//From Account Balance Value After Transaction
	TVA					big.Int						//To Account Balance Value After Transaction
	Processor			address
	State				string						//Sent | Confirming | Confirmed | Balance To Small
	Fee					big.Int
	Date				time.Time				 	 //UTC time transaction took place tXID Date and time will use server timezone	
	Signature			SignedTx					 //Signature of Sender
	TxConfirmation		map[unit64][]SignedTx	 	 //Map is Hash of Transaction and the associated array of verfied signatures of validators	
	TxHash				uint64

}

type Transactions struct{ 
	ChainID		  uint
	Transactions  map[uint64]Transaction  // Map of all Transactions for the Year by TxHash

}

type TransactionPool struct{
   Txs			[]Transaction
   NextNumber	uint64
}


func createTxPool() TransactionPool{
	txs := []Transaction
	txPool := TransactionPool{txs, 1}
}
func (txPool *TransactionPool ) addTxs(transaction Transaction){
		txPool.Txs.append(transaction)
		txPool.NextNumber +=1 
}
func getBalance(addr Common.Address){
	account := lookUpAccount(addr)
	return account.Balance
}
func sendTransaction(to Common.Address, from Common.Address, amount big.Int, key *ecdsa.PrivateKey, auth string){
		toBalance :=getBalance(to)
		fromBalance :=getBalance(from)
		if fromBalance.Cmp(amount)>0{
		
			return "Transaction sent " + string(tx.TxHash)
		}else{
			return "Error not enough FGEs have " + fromBalance.String() +" Need "+amount.String()
			
		}
}

func SaveTransactionToDisk(tx Transaction){


}

func LoadTransactionFromDisk(hash uint64){


}