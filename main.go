package main

import(
	"bytes"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	//"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/node"
	"github.com/fgeth/fg/transaction"
)
var (

	wg 		 sync.WaitGroup

)

func main(){
	
	common.MyNode = NewNode()
	common.MyNode.SaveNode()
	if len(os.Args)>1{
	switch os.Args[1]{
	case "Gen":
		fmt.Println("Genesis Block")
		//common.ImportBlocks()
		//common.ImportTxs()
		common.SignGenesisBlocks() 
	case "Node":
		time.Sleep(time.Second * 60)
		os.Exit(0)

	}
	}else{
		go register()
		//common.ImportBlocks()
		//common.ImportTxs()
		//common.GetBlocks()
		//common.GetTxs()
	}
	wg.Add(1)
	server()
	go fg()
	go CloseHandler()
	wg.Wait()
	
}
func server(){
	wg.Add(1)
	defer wg.Done()
	r := mux.NewRouter()
    
	//r.HandleFunc("/", home).Methods("GET")
	//r.HandleFunc("/getBlocks", sendBlocks).Methods("GET")
	//r.HandleFunc("/getNodes", sendNodes).Methods("GET")
	//r.HandleFunc("/getTxs", sendTxs).Methods("GET")
	r.HandleFunc("/sendTx", sendNewTransaction).Methods("POST")
	//r.HandleFunc("/Tx", CreateNewTransaction).Methods("POST")
	//r.HandleFunc("/block", createNewBlock).Methods("POST")
	//r.HandleFunc("/newNode", newNode).Methods("POST")
	//r.HandleFunc("/blockTxs", processTxs).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Listening on port 42069")
if err := http.ListenAndServe(":42069", nil); err != nil {
	    	log.Fatal(err)
	   }
	   
}

//Function to check that new blocks are being made and if not start the process
func fg(){
	wg.Add(1)
	defer wg.Done()
	for {
		CheckBlockNumber := common.BlockNumber	
		time.Sleep(time.Second * 60)
		if CheckBlockNumber == common.BlockNumber{
			common.BlockFailed(common.BlockNumber)
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
//TODO Fix this
func register() bool{
	haveNode := false
	ImportActiveNodes()
	if len(common.ActiveNodes) >0{
		haveNode = true	
	}else{
		//common.MyNode.GetNodes()
		if len(common.ActiveNodes) >0{
			haveNode = true	
		}
	}	
	if haveNode {
	for x:=0; x<len(common.ActiveNodes); x+=1{
			
			//common.MyNode.RegisterNode(common.ActiveNodes[x])
		}
	}
	return haveNode
}
func SaveFiles(){
//TODO Implement way to save current BlockChain State to File
	SaveActiveNodes()
	
}


func SaveActiveNodes(){

	file, _ := json.MarshalIndent(common.ActiveNodes, "", " ")
 
	_ = ioutil.WriteFile("ActiveNodes.json", file, 0644)

}
func ImportActiveNodes(){
	file, _ := ioutil.ReadFile("ActiveNodes.json")
	_ = json.Unmarshal([]byte(file), &common.ActiveNodes)
}

//TODO fix this
func sendNewTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var tx transaction.Transaction
    json.Unmarshal(reqBody, &tx)
	if common.VaildTransaction(tx){
		if common.MyNode.Writer{
			common.BTxHash = append(common.BTxHash, tx.TxHash)
		
		}else{	
			//common.SendTransaction(tx)
		}
	
		json.NewEncoder(w).Encode(tx)
	}
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var tx transaction.Transaction
    json.Unmarshal(reqBody, &tx)
}

func buildTransaction(){
	
}

func NewNode() *node.Node{
	var node *node.Node
	return node
}
func newBlock(blockNumber uint64){
	common.CreateBlock( blockNumber )
}

func verifyTx(OTP string, Tx transaction.Transaction) transaction.Transaction{
	var inValidTx transaction.Transaction 
	txCreditAmt :=big.NewInt(0)
	txDebitAmt :=big.NewInt(0)
	txCAmt :=big.NewInt(0)
	txDAmt :=big.NewInt(0)
	txAcc :=big.NewInt(0)
	for x:=0; x< len(Tx.Credit); x+=1{
		txCreditAmt.Add(txCreditAmt, Tx.Credit[x].Amount)
	}
	for x:=0; x< len(Tx.Debit); x+=1{
		txDebitAmt.Add(txDebitAmt, Tx.Debit[x].Amount)
	}
	txInt := Tx.CalcInterest()
	txFee := Tx.CalcFee()
	txCAmt.Add(txInt, txCreditAmt)
	txDAmt.Add(txFee, txDebitAmt)
	Tx.Change.Amount.Sub(txCAmt, txDAmt)
	txAcc.Add(txDAmt,Tx.Change.Amount)
	txAcc.Sub(txAcc,txCAmt)
	if (txAcc.Cmp(big.NewInt(0))==0){

		if bytes.Compare(Tx.TxHash, Tx.HashTx()) ==0 {
			if Tx.Payout == false{
				return Tx
				}
			}
	}
	return inValidTx
}

