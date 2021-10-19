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
	 "github.com/fgeth/fg/crypto"
)
var (
	
)

type BaseTransaction struct {
	ChainYear			uint32						//Chain Transaction belongs to	
	BlockNumber			uint64						//Block Number Transaction was created in
	Time				time.Time					//Time Transaction was Created time.now()
	Amount				*big.Int					//Amount in FGs
	TxHash				crypto.Hash					//Hash of  ChainYear, Time, and amount plus OTP if Debit Transaction
	Spent				crypto.Hash					//Hash of Transaction were the Debit Balance of this Transaction was spent
	TxId				crypto.Hash					//Hash the Transaction that this Base Transaction is a part of as a Debit Transaction
	
	
}

type Transaction struct {
	TxHash				crypto.Hash					//Hash of Credit, Debit, and Change hashes plus OTP of sender
    Debit				[]BaseTransaction
	Change				  BaseTransaction				//Debit Transaction to give any change due to sender
	Credit				[]BaseTransaction
	OTP					string						//One Time Password which is the Public Key as a string Used for this Transaction if present transaction value has been spent
	R					big.Int						//Part one of Signature of sender when they sign the transaction Hash
	S					big.Int						//Part two of Signature
	Payout				bool
}
func (tx *BaseTransaction) SaveTx(){
    dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
 
	path :=filepath.Join(dirname, "fg", "btx")
	 
	folderInfo, err := os.Stat(path)
	if folderInfo.Name() !="" {
			fmt.Println("")
	}
    if os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(dirname, "fg"), 0755)
		fmt.Println(err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
	uintA, uintB, uintC, uintD := crypto.HashToUint64(tx.TxHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	fileName := filepath.Join(path,theHash)
	fmt.Println(fileName)
	file, _ := json.MarshalIndent(tx, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}
func ImportBaseTx(txHash crypto.Hash) BaseTransaction{
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
	uintA, uintB, uintC, uintD := crypto.HashToUint64(txHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	path :=filepath.Join(dirname, "fg", "btx", theHash )
	file, _ := ioutil.ReadFile(path)
	var tx BaseTransaction
	_ = json.Unmarshal([]byte(file), &tx)
	
	return tx
}

func (tx *Transaction) SaveTx(){
    dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
 
	path :=filepath.Join(dirname, "fg", "tx")
	 
	folderInfo, err := os.Stat(path)
	if folderInfo.Name() !="" {
			fmt.Println("")
	}
    if os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(dirname, "fg"), 0755)
		fmt.Println(err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
	uintA, uintB, uintC, uintD := crypto.HashToUint64(tx.TxHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	fileName := filepath.Join(path,theHash)
	fmt.Println(fileName)
	file, _ := json.MarshalIndent(tx, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}
func ImportTx(txHash crypto.Hash) Transaction{
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
	uintA, uintB, uintC, uintD := crypto.HashToUint64(txHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	path :=filepath.Join(dirname, "fg", "tx", theHash )
	file, _ := ioutil.ReadFile(path)
	var tx Transaction
	_ = json.Unmarshal([]byte(file), &tx)
	
	return tx
}


func (tx Transaction) CalcFee()big.Int{
percentage := big.NewInt(100)
txFee :=big.NewInt(0)
for x:=0; x< len(tx.Debit); x+=1{
	txFee.Add(txFee, new(big.Int).Div(tx.Debit[x].Amount, percentage))
}
return txFee
}
func (tx Transaction) CalcInterest() big.Int{
    interest  :=big.NewInt(0)
	txInterest:=big.NewInt(0)
	for x:=0; x< len(tx.Credit); x+=1{
		months := int64(time.Now().Sub(tx.Credit[x].Time)/(720*time.Hours))
		if months >0{
			m := big.NewInt(months)
			interest.Add(interest.Mul(new(big.Int).Div(tx.Credit[x].Amount, percentage), big.NewInt(2))
			txInterest.Add(txInterest, interest.Mul(interest, m))
		}
	}
	return txInterest
 }

func (Tx Transaction) Credits() big.Int{
txAmount :=big.NewInt(0)
for x:=0; x< len(tx.Credit); x+=1{
	txAmount.Add(txAmount, Tx.Credit[x].Amount)
}
	return txAmount
}

func (Tx Transaction) Debits() big.Int{
txAmount :=big.NewInt(0)
for x:=0; x< len(tx.Debit); x+=1{
	txAmount.Add(txAmount, Tx.Debit[x].Amount)
}
	return txAmount
}



func(Tx BaseTransaction) HashBaseTx(pubKey string ) crypto.Hash{
	kh crypto.NewKeccakState()
	txData := string(Tx.ChainYear) + Tx.Time.String() + Tx.Amount.String() + pubKey
	return crypto.HashData(kh, []byte(TxData)){

}

func(Tx Transaction) HashTx( ) crypto.Hash{
	kh crypto.NewKeccakState()
	for x:=0; x< len(Tx.Credit); x+=1{
		txData := string(Tx.Credit[x].TxHash) 
		
	}
	for x:=0; x< len(Tx.Debit); x+=1{
		txData = txData + string(Tx.Debit[x].TxHash) 
		
	}
	txData = txData + string(Tx.Change.TxHash)
	txData = txData + Tx.OTP
	return crypto.HashData(kh, []byte(txData)){

}

func(Tx Transaction) VerifySig() bool{
	kh crypto.NewKeccakState()
		
	publicKey := DecodePubKey(Tx.OTP)
	if Tx.TxHash == Tx.HashTx(){
		return crypto.verify(Tx.TxHash, Tx.R, Tx.S, publicKey) 

	}else{
		return FALSE
	}
	
}