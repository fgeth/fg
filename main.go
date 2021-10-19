package main

import(

	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/node"
	"github.com/fgeth/fg/transaction"
)
var (
	node	 *Node
	wg 		 sync.WaitGroup

)

func main(){
	
	node = NewNode()
	SaveNode()
	switch os.Args[1]{
	case "Gen":
		ImportBlocks()
		ImportTxs()
		node.SignGenesisBlocks() 
	case "Node":
		time.Sleep(time.Second * 60)
		os.Exit(0)
	default:
		go register()
		ImportBlocks()
		ImportTxs()
		node.GetBlocks()
		node.GetTxs()
	}
	wg.Add(1)
	go server()
	go fg()
	go CloseHandler()
	wg.Wait()
	
}
func server(){
	wg.Add(1)
	defer wg.Done()
	r := mux.NewRouter()
    
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/getBlocks", sendBlocks).Methods("GET")
	r.HandleFunc("/getNodes", sendNodes).Methods("GET")
	r.HandleFunc("/getTxs", sendTxs).Methods("GET")
	r.HandleFunc("/sendTx", sendNewTransaction).Methods("POST")
	r.HandleFunc("/Tx", CreateNewTransaction).Methods("POST")
	r.HandleFunc("/block", createNewBlock).Methods("POST")
	r.HandleFunc("/newNode", newNode).Methods("POST")
	r.HandleFunc("/blockTxs", processTxs).Methods("POST")
	http.Handle("/", r)
if err := http.ListenAndServe(":42069", nil); err != nil {
	    	log.Fatal(err)
	   }
	   
}

//Function to check that new blocks are being made and if not start the process
func fg(){
	wg.Add(1)
	defer wg.Done()
	for {
		CheckBlockNumber := BlockNumber	
		time.Sleep(time.Second * 60)
		if CheckBlockNumber == BlockNumber{
			node.BlockFailed(BlockNumber)
		}
	}

}

func CloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal Closing Program")
		SaveFiles()
		os.Exit(0)
	}()
}

func register() bool{
	haveNode := FALSE
	ImportActiveNodes()
	if len(ActiveNodes) >0{
		haveNode = TRUE	
	}else{
		node.GetNodes()
		if len(ActiveNodes) >0{
			haveNode = TRUE	
		}
	}	
	if haveNode {
	for x:=0; x<len(ActiveNodes); x+=1{
			node.RegisterNode(ActiveNodes[x])
		}
	}
}
func SaveFiles(){
//TODO Implement way to save current BlockChain State to File
	SaveActiveNodes()
	
}


func SaveActiveNodes(){

	file, _ := json.MarshalIndent(ActiveNodes, "", " ")
 
	_ = ioutil.WriteFile("ActiveNodes.json", file, 0644)

}
func ImportActiveNodes(){
	file, _ := ioutil.ReadFile("ActiveNodes.json")
	_ = json.Unmarshal([]byte(file), &ActiveNodes)
}

func sendNewTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var tx Transaction
    json.Unmarshal(reqBody, &tx)
	if tx.VaildTransaction()(){
		if node.Writer{
			BTx = append(BTx, tx.TxHash)
		
		}else{	
			node.SendTransaction(tx)
		}
	
		json.NewEncoder(w).Encode(tx)
	}
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var tx Transaction
    json.Unmarshal(reqBody, &tx)
}

func (node *Node) buildTransaction(){
	
}


func verifyTx(OTP string, Tx Transaction) Transaction{
	
	txCreditAmt :=big.NewInt(0)
	txDebitAmt :=big.NewInt(0)
	txCAmt :=big.NewInt(0)
	txDAmt :=big.NewInt(0)
	txAcc :=big.NewInt(0)
	for x:=0; x< len(Tx.Credits); x+=1{
		txCreditAmt := txCreditAmt.Add(txCreditAmt, Tx.Credits[x].Amount)
	}
	for x:=0; x< len(Tx.Debits); x+=1{
		txDebitAmt := txdebitAmt.Add(txdebitAmt, Tx.Debits[x].Amount)
	}
	txInt := Tx.CalcInterest()
	txFee := Tx.CalcFee()
	txCAmt.Add(txInt, txCreditAmt)
	txDAmt.Add(txFee, txDebitAmt)
	txChange.Sub(txCAmt, txDAmt)
	txAcc.add(txDAmt,Tx.Change.Amount)
	txAcc.Sub(txAcc,txCAmt)
	if (txAcc == 0){

		if Tx.TxHash == Tx.HashTx(){
			if Tx.Payout == FALSE{
				return Tx
				}
			}
	}
	
}


func (tx Transaction) VaildTransaction() bool{
txFee := tx.CalcFee()

if FGValue =>100{
	txInterest := tx.CalcInterest()
}else{
	txInterest :=0
}

if Tx.OTP !=""{
		if Tx.VerifySig(){
			if Tx.Credits().Add(Tx.Credits,txInterest) == txFee.Add(txFee, Tx.Debits().Add(Tx.Debits(), Tx.Change)){
			for x:=0; x< len(tx.Credit); x+=1{
				cTx := ImportBaseTx(tx.Credit[x].TxHash)
				if cTx.OTP ==""{
					cTx.OTP=tx.Credit[x].OTP
					cTx.SaveTx()
					
				}else{
					return FALSE		//Double Spend
				}
			}
			for x:=0; x< len(tx.Debit); x+=1{
				tx.Debit[x].SaveTx()

			}
				tx.Change.SaveTx()
			}else{
					return FALSE
			}
		}else{
			return FALSE
		}
}else{
	if Tx.Credit.Payout == TRUE{
		if tx.VerifySig(){
			block := ImportBlock(tx.Credit.ChainYear , tx.Credit.BlockNumber)
			for x:=0; x<len(block.Writers); x +=1{
				if block.Writers[x] == tx.OTP{
					if tx.Debit.Amount == block.NodePayout
						return TRUE
					if tx.Debit.Amount == block.WriterPayout
				}
			}
		}
	
	}
	return FALSE
}
}


func (node *Node) CreatePayoutTransaction(amt big.Int, pubKey string, blockNumber uint64) Transaction{

	var C []BaseTransaction
	var D []BaseTransaction
	var Debit BaseTransaction
	var Credit BaseTransaction
	var TX Transaction
	Debit.ChainYear = ChainYear
	Debit.BlockNumber = blockNumber
	Debit.Time = time.Now()
	Debit.Amount = amt
	Debit.TxHash = Debit.HashBaseTx(pubKey)
	Credit.ChainYear = ChainYear
	Credit.BlockNumber = blockNumber
	Credit.Time = time.Now()
	Credit.TxHash = Credit.HashBaseTx(pubKey)
	Credit.Amount = amt
	Tx.Change.TxHash = Tx.Change.HashBaseTx(pubKey)
	Tx.Debit =append(Tx.Debit, Debit)
	Tx.Credit = append(Tx.Credit, Credit)
	Tx.OTP = node.Id
	Tx.TxHash = Tx.HashTx(Tx.OTP)
	Tx.R, Tx.S err := crypto.sign(Tx.Hash, node.PrvKey)
	Tx.Payout = TRUE
	return Tx
	
}
func (node *Node) sellItem(item Item) {
	prvKey := GenerateRSAKey()
	node.Comms.RsaPrvKeys[prvKey.PublicKey] = prvKey
	item.Seller = prvKey.PublicKey
	node.Items.Selling.Item = append(node.Items.Selling.Item, item)

	
}
