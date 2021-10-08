package transaction

import (
    "crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
    "encoding/pem"
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
	TxId				string						//time.Now.UTC.String()::NodeID::TransactionNumber
	BlockNumber			unit64						//Block Number transaction was confirmed
	From				address
	To					address
	Value				big.Int
	TxNumber			uint64						//The senders total number of sent transactions to this point including this transaction
	Processor			address
	State				int							//-1 Rejected | 1 Sent | 2 Confirming | 3 Completed | 4 Confirmed 
	Fee					big.Int
	Date				time.Time				 	 //UTC time transaction took place TxId Date and time will use server timezone	
	Signature			Signer						 //Signature of Sender
	TxHash				uint64
	Confirmations		[]SignedTx					//Array of Node signatures that have comfirmed this Transaction
	Challenged			bool						//Used when a Transaction is Challenged
	

}

type Transactions struct{ 
	ChainID		  uint
	Transactions  map[uint64]Transaction  // Map of all Transactions for the Year by TxHash

}

type TransactionPool struct{
   Txs			[]Transaction
   NextNumber	uint64
}

type SignedTx struct {
    Accept			bool				//If node accepts this transaction or rejects the transaction
	R				big.Int
	S				big.Int
	Node			uintptr  			//Able to look up Node and get its public key
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




func SaveTransactionToDisk(tx Transaction){


}

func LoadTransactionFromDisk(hash uint64){


}
func Address2Key(address string) ( *ecdsa.PublicKey) {
    blockPub, _ := pem.Decode([]byte(address))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return publicKey
}
func (Tx *Transaction) VerifySender(){
	
}
func (Tx *Transaction) GetUnsignedTransaction()  *Transaction{

	checkTrans :=Transaction{}
	checkTrans.From := Tx.From
	checkTrans.To := Tx.To
	checkTrans.Value := Tx.Value
	checkTrans.TxNumber := Tx.TxNumber
	return checkTrans
}

func (Tx *Transaction) VerifyTransactionHash() bool{
	
	kh := NewKeccakState()
	
	data, _ := json.Marshal(Tx.GetUnsignedTransaction())
	h := HashData(kh , []byte(data))
	return Tx.TxHash == h
}




