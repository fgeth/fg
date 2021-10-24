package transaction

import (
	"bytes"
    "fmt"
	"io/ioutil"
	"math/big"
	"time"
	"encoding/json"
	 "os"
	 "path/filepath"
	 "strconv"
	 "github.com/fgeth/fg/crypto"
)
var (
	
)

type BaseTransaction struct {
	ChainYear			uint64						//Chain Transaction belongs to	
	BlockNumber			uint64						//Block Number Transaction was created in
	Time				time.Time					//Time Transaction was Created time.now()
	Amount				*big.Int					//Amount in FGs
	TxHash				string						//Hash of  ChainYear, Time, and amount plus OTP if Debit Transaction
	Spent				string						//Hash of Transaction were the Debit Balance of this Transaction was spent
	TxId				string						//Hash the Transaction that this Base Transaction is a part of as a Debit Transaction
	OTP					string						//Public Key as string
	R					*big.Int						//Part one of Signature of sender when they sign the transaction Hash
	S					*big.Int						//Part two of Signature
	
}

type Transaction struct {
	TxHash				string				//Hash of Credit, Debit, and Change hashes plus OTP of sender
    Debit				BaseTransaction
	Change				BaseTransaction					//Debit Transaction to give any change due to sender
	Credit				[]BaseTransaction
	Payout				bool
}

type BaseTxData struct{
    ChainYear			uint64						//Chain Transaction belongs to	
	BlockNumber			uint64						//Block Number Transaction was created in
	Time				time.Time					//Time Transaction was Created time.now()
	Amount				*big.Int					//Amount in FGs
	OTP					string						//One Time Password which is the Public Key as a string Used for this Transaction if present transaction value has been spent

}

type TxData struct{
    Debit				BaseTransaction
	Change				BaseTransaction				//Debit Transaction to give any change due to sender
	Credit				[]BaseTransaction
	

}
func (tx *BaseTransaction) SaveTx(dirname string){
    //dirname, err := os.UserHomeDir()
    //if err != nil {
    //    fmt.Println( err )
    //}
 
	path :=filepath.Join(dirname, "btx")
	 
	folderInfo, err := os.Stat(path)
	if folderInfo.Name() !="" {
			fmt.Println("")
	}
    if os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		fmt.Println(err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
	uintA, uintB, uintC, uintD := crypto.B32HashToUint64([]byte(tx.TxHash))
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
func ImportBaseTx(txHash []byte, dirname string) BaseTransaction{
	//dirname, err := os.UserHomeDir()
    //if err != nil {
    //    fmt.Println( err )
    //}
	uintA, uintB, uintC, uintD := crypto.B32HashToUint64(txHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	path :=filepath.Join(dirname, "btx", theHash )
	file, _ := ioutil.ReadFile(path)
	var tx BaseTransaction
	_ = json.Unmarshal([]byte(file), &tx)
	
	return tx
}

func (tx *Transaction) SaveTx(dirname string){
    //dirname, err := os.UserHomeDir()
    //if err != nil {
    //    fmt.Println( err )
    //}


	path :=filepath.Join(dirname, "tx")
	 
	folderInfo, err := os.Stat(path)
    if os.IsNotExist(err) {

		err := os.Mkdir(dirname, 0755)
		fmt.Println(err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
		}else{
			fmt.Println(folderInfo)
		}
		fmt.Println("TxHash:",len(tx.TxHash))
		uintA, uintB, uintC, uintD := crypto.B32HashToUint64([]byte(tx.TxHash))
		h1 := strconv.FormatUint(uintA, 10)
		h2 := strconv.FormatUint(uintB, 10)
		h3 := strconv.FormatUint(uintC, 10)
		h4 := strconv.FormatUint(uintD, 10)
		theHash := h1 + h2 +h3 +h4
		fileName := filepath.Join(path,theHash)
		fmt.Println(fileName)
		file, _ := json.MarshalIndent(tx, "", " ")
	 
		_ = ioutil.WriteFile(fileName, file, 0644)
		
		
			uintA, uintB, uintC, uintD = crypto.B32HashToUint64([]byte(tx.Debit.TxHash))
			h1 = strconv.FormatUint(uintA, 10)
			h2 = strconv.FormatUint(uintB, 10)
			h3 = strconv.FormatUint(uintC, 10)
			h4 = strconv.FormatUint(uintD, 10)
			theHash = h1 + h2 +h3 +h4
			fileName = filepath.Join(path,theHash)
			fmt.Println(fileName)
			file, _ = json.MarshalIndent(tx.Debit, "", " ")
		
		_ = ioutil.WriteFile(fileName, file, 0644)
	
		for x:=0; x < len(tx.Credit); x +=1{
			uintA, uintB, uintC, uintD = crypto.B32HashToUint64([]byte(tx.Credit[x].TxHash))
			h1 = strconv.FormatUint(uintA, 10)
			h2 = strconv.FormatUint(uintB, 10)
			h3 = strconv.FormatUint(uintC, 10)
			h4 = strconv.FormatUint(uintD, 10)
			theHash = h1 + h2 +h3 +h4
			fileName = filepath.Join(path,theHash)
			fmt.Println(fileName)
			file, _ = json.MarshalIndent(tx.Credit[x], "", " ")
		 
			_ = ioutil.WriteFile(fileName, file, 0644)
		}
		
		
			uintA, uintB, uintC, uintD = crypto.B32HashToUint64([]byte(tx.Change.TxHash))
		h1 = strconv.FormatUint(uintA, 10)
		h2 = strconv.FormatUint(uintB, 10)
		h3 = strconv.FormatUint(uintC, 10)
		h4 = strconv.FormatUint(uintD, 10)
		theHash = h1 + h2 +h3 +h4
		fileName = filepath.Join(path,theHash)
		fmt.Println(fileName)
		file, _ = json.MarshalIndent(tx.Change, "", " ")
	 
		_ = ioutil.WriteFile(fileName, file, 0644)

}
func ImportTx(txHash crypto.Hash, dirname string) Transaction{
	//dirname, err := os.UserHomeDir()
    //if err != nil {
    //    fmt.Println( err )
    //}
	uintA, uintB, uintC, uintD := crypto.HashToUint64(txHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	path :=filepath.Join(dirname, "tx", theHash )
	file, _ := ioutil.ReadFile(path)
	var tx Transaction
	_ = json.Unmarshal([]byte(file), &tx)
	
	return tx
}


func (tx Transaction) CalcFee() *big.Int{
	percentage := big.NewInt(500)
	txFee :=big.NewInt(0)
	txFee.Add(txFee, new(big.Int).Div(tx.Debit.Amount, percentage))
	return txFee
}
func (tx Transaction) CalcInterest() *big.Int{
    interest  :=big.NewInt(0)
	txInterest:=big.NewInt(0)
	percentage := big.NewInt(100)
	for x:=0; x< len(tx.Credit); x+=1{
		months := int64(time.Now().Sub(tx.Credit[x].Time)/(720*time.Hour))
		if months >0{
			m := big.NewInt(months)
			q:= new(big.Int).Div(tx.Credit[x].Amount, percentage)
			interest.Mul(q, big.NewInt(2))
			interest.Add(interest, q)
			interest.Mul(interest, m)
			txInterest.Add(txInterest, interest )
		}
	}
	return txInterest
 }

func (Tx Transaction) Credits() *big.Int{
txAmount :=big.NewInt(0)
for x:=0; x< len(Tx.Credit); x+=1{
	txAmount.Add(txAmount, Tx.Credit[x].Amount)
}
	return txAmount
}

func (Tx Transaction) Debits() *big.Int{
txAmount :=big.NewInt(0)

	txAmount.Add(txAmount, Tx.Debit.Amount)
	return txAmount
}



func(Tx BaseTransaction) HashBaseTx(pubKey string ) string{
	//kh :=crypto.NewKeccakState()
	txData := Tx.BaseTxData()
	txData.OTP = pubKey
	json,_ := json.Marshal(txData)

	return crypto.HashTx([]byte(json))

}

func(Tx Transaction) HashTx( ) string{
	//kh :=crypto.NewKeccakState()
	txData := Tx.TxData()
	json,_ := json.Marshal(txData)
	return crypto.HashTx([]byte(json))

}


func (Tx BaseTransaction) BaseTxData() BaseTxData{
    var txData BaseTxData
	
	txData.ChainYear 		= Tx.ChainYear
	txData.BlockNumber 		= Tx.BlockNumber
	txData.Time				= Tx.Time
	txData.Amount			= Tx.Amount
	 
	
	return txData
}

func (Tx Transaction) TxData() TxData{
    var txData TxData
	
	txData.Credit = Tx.Credit
	fmt.Println("Tx Credit :", txData.Credit)
	txData.Debit = Tx.Debit
	txData.Change = Tx.Change
	
	return txData
}

func(Tx BaseTransaction) VerifySig() bool{
			
	publicKey := crypto.DecodePubKey(Tx.OTP)
	if bytes.Compare([]byte(Tx.TxHash), []byte(Tx.HashBaseTx(Tx.OTP))) ==0 {
		return crypto.TxVerify([]byte(Tx.Spent), Tx.R, Tx.S, publicKey) 

	}else{
		return false
	}
	
}