package main

import(
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"flag"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	//"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/node"
	"github.com/fgeth/fg/transaction"
	"github.com/fgeth/bine/tor"


)
var (

	wg 		 sync.WaitGroup
	path	  string
	port 	  string
	torPort	  string 		//Port Tor is running on
	Gen		  *bool
	
	
)

func init() {
			flag.StringVar(&port, "port", "42069", "Default Port")
			flag.StringVar(&path, "path", "/var/fg", "Data Directory")
			flag.StringVar(&torPort, "torPort", "9051", "Default Port")
	Gen = 	flag.Bool("gen", false, "Continue with existing chain")
	
}

func main(){	
	flag.Parse()
	
	common.MyNode = node.ImportNode(path)
	
	if common.MyNode.Id ==""{
		common.MyNode = NewNode()

	}else{
		if common.MyNode.PKStr !=""{
			common.MyNode.PubKey = crypto.DecodePubKey(common.MyNode.PKStr)
		}
		if common.MyNode.PRKStr !=""{
			common.MyNode.PrvKey = crypto.DecodePrv(common.MyNode.PRKStr)
		}
	}
	Trusted()
	if port !="42069"{
		common.MyNode.Port = ":"+port
		}
	if torPort !="9051"{
		common.MyNode.Tor = ":"+torPort
		}
	
	directory()
	wallet()
	common.MyNode.SaveNode(common.MyNode.Path)
	fmt.Println("Node Id is :" , common.MyNode.Id)
	//fmt.Println("Node Ip is :" , common.MyNode.Ip)
	fmt.Println("Node Path is :" , common.MyNode.Path)
	test()
	
	if *Gen{
		fmt.Println("Genesis Block")
		common.FGValue = .01
		//common.ImportBlocks()
		//common.ImportTxs()
		common.SignGenesisBlocks() 
	}else{
		go register()
		//common.ImportBlocks()
		//common.ImportTxs()
		//common.GetBlocks()
		//common.GetTxs()
	}
	wg.Add(1)
	go server()
	torServer()
	
	
	//TorService()
	//time.Sleep(time.Second * 120)
	
	go fg()
	go CloseHandler()
	wg.Wait()
	
}

func Trusted(){
	fg1 :="-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE7mqd2vQKqyk/Xy171L/0YFnqT4Uy\nfXJeUmK41oNakzblaAbOTjnrVwt3LRt1tdtl4Um+gld87Vz860GcF4os+w==\n-----END PUBLIC KEY-----\n"
	pubKey := crypto.DecodePubKey(fg1)
	common.Trusted = append(common.Trusted, pubKey)
	add := crypto.GetAddress(pubKey)
	//key := crypto.UnSetBytes(add)
	fmt.Println("PubKey Address :", add)
	addHash := crypto.HashTx([]byte(add))
	R, S  := crypto.TxSign([]byte(addHash), common.MyNode.PrvKey)
	trusted := crypto.TxVerify([]byte(addHash), R, S, pubKey)
	fmt.Println("This is a trusted Node :", trusted)
	//trusted = crypto.TxVerify([]byte(addHash), R, S, key)
	//fmt.Println("This is a trusted Node :", trusted)
	//0xe56806Ce4e39Eb570b772d1B75B5dB65e149be82577ceD5cfceB419c178A2cFb
							  //0x75b5Db65e149be82577CED5CfCeb419c178a2cFb
}

func test(){
	BlockReward:= big.NewInt(0)
	BlockReward.SetString("10000000000000000000", 10)
	bn :=common.BlockNumber + uint64(1)
	k,_:=crypto.GenerateKey()
	prvK, pubK := crypto.Encode(k,&k.PublicKey)
	fmt.Println("Pvt Key:", prvK)
	var keys  []*ecdsa.PrivateKey
	keys = append(keys, k)
	credit := common.CreateDebitTxs(BlockReward, pubK, bn)
	var credits []transaction.BaseTransaction
	credits = append(credits, credit)
	tx1 :=common.CreateTransaction(BlockReward, credits, pubK,pubK, bn, keys )
	fmt.Println("Tx ", tx1)
	tx1.SaveTx(common.MyNode.Path)
	add := crypto.BytesToAddress([]byte(tx1.TxHash))
	fmt.Println("Address :", add)
	common.FGValue = .01
	common.Wallet.FGs = common.Wei2FG(BlockReward)
	common.Wallet.Wei = BlockReward
	common.Wallet.Dollars = common.FG2USD(BlockReward)

}
func directory(){

	fmt.Println("Creating Directories :", common.MyNode.Path) 
	_, err := os.Stat(common.MyNode.Path)

    if os.IsNotExist(err) {
		err = os.Mkdir(common.MyNode.Path, 0755)
		fmt.Println(err)
		
    }
	path =filepath.Join(common.MyNode.Path, "node")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		fmt.Println(err)
		
    }
	
	path =filepath.Join(common.MyNode.Path, "block")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		fmt.Println(err)
		
    }
	
	path =filepath.Join(common.MyNode.Path, "tx")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		fmt.Println(err)
		
    }
	
	path =filepath.Join(common.MyNode.Path, "store")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		fmt.Println(err)
		
    }

}
func wallet(){
common.Wallet.Items.Item = map[string]item.Item{}
common.Wallet.Items.Keys = map[string][]*ecdsa.PrivateKey{}

}
func server(){
	wg.Add(1)
	defer wg.Done()
	
    r := mux.NewRouter()
	//r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/getBlocks", sendBlocks).Methods("GET")
	//r.HandleFunc("/getNodes", sendNodes).Methods("GET")
	//r.HandleFunc("/getTxs", sendTxs).Methods("GET")
	r.HandleFunc("/getWallet", GetWallet).Methods("GET")
	
	r.HandleFunc("/sendTx", sendNewTransaction).Methods("POST")
	//r.HandleFunc("/Tx", CreateNewTransaction).Methods("POST")
	r.HandleFunc("/block", createNewBlock).Methods("POST")
	r.HandleFunc("/addItem", createNewItem).Methods("POST")
	//r.HandleFunc("/newNode", newNode).Methods("POST")
	//r.HandleFunc("/blockTxs", processTxs).Methods("POST")
	
	
	
	staticFileDirectory := http.Dir(path)
	
	staticFileHandler := http.StripPrefix("/store/", http.FileServer(staticFileDirectory))
	
	r.PathPrefix("/store/").Handler(staticFileHandler).Methods("GET")
	
	fmt.Println("Listening on port :", 80)
	server := &http.Server{
    Addr:    ":80",
    Handler: r,
		}
if err := server.ListenAndServe(); err != nil {
	    	log.Fatal(err)
	   }
	   
}

func torServer() error {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(nil, nil,common.MyNode.Path )
	r := mux.NewRouter()
	if err != nil {
		return err
	}
	defer t.Close()
	// Add a handler
	http.Handle("/", r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Dark World!"))
	})
	r.HandleFunc("/sendTx", sendNewTransaction).Methods("POST")
	r.HandleFunc("/block", createNewBlock).Methods("POST")
	r.HandleFunc("/addItem", createNewItem).Methods("POST")
	
	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer listenCancel()
	// Create an onion service to listen on 42069 but show as 420
	port, _ := strconv.Atoi(common.MyNode.Port)
	onion, err := t.Listen(listenCtx, &tor.ListenConf{LocalPort: port , RemotePorts: []int{80}, Version3: true})
	if err != nil {
		return err
	}
	defer onion.Close()
	// Serve on HTTP
	fmt.Printf("Listening on port :", common.MyNode.Port)
	fmt.Printf("Open Tor browser and navigate to http://%v.onion\n", onion.ID)
	common.MyNode.OA = onion.ID
	return http.Serve(onion, nil)
	
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

func createNewBlock(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var block block.Block
    json.Unmarshal(reqBody, &block)
	if common.VerifyBlock(&block){
		block.SaveBlock(common.MyNode.Path)
		common.BlockNumber = block.BlockNumber
		
		json.NewEncoder(w).Encode(block)
	}
}

func sendBlocks(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var block2 block.Block
    json.Unmarshal(reqBody, &block2)
	Block := block.ImportBlock(block2.ChainYear, block2.BlockNumber, common.MyNode.Path)
	json.NewEncoder(w).Encode(Block)
}


func createNewItem(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var theItem item.Item
	aHash := crypto.HashTx(reqBody)
	fmt.Println("aHash", aHash)
	//fmt.Println("Item", reqBody)
    json.Unmarshal(reqBody, &theItem)
	// _, _, _, uintD :=crypto.B32HashToUint64(aHash)
	// item.Id =  strconv.FormatUint(uintD, 10)
	theItem.Id = aHash
	fmt.Println("Amount", theItem.Amount)
	
	theItem.Tx.Tx = map[string][]transaction.BaseTransaction{theItem.Id: createDebitTx(theItem.Amount, theItem)}
	prvKey := crypto.GenerateRSAKey()
	pubKey := prvKey.PublicKey
	theItem.Seller = pubKey
	path := filepath.Join(common.MyNode.Path, "Keys", theItem.Id) 
	crypto.StoreRSAKey(prvKey, "Password", path)
	theItem.SaveItem(common.MyNode.Path)
	
	common.Wallet.Items.Item[theItem.Id] = theItem
	fmt.Println("New Item")
	fmt.Println("The Item: ", theItem)
	fmt.Println("Num Debit Tx ", len(theItem.Tx.Tx[theItem.Id]))
	for x :=0; x < len(theItem.Tx.Tx[theItem.Id]); x +=1{
			fmt.Println("Debit Amount", theItem.Tx.Tx[theItem.Id][x].Amount)
			fmt.Println("Debit OTP", theItem.Tx.Tx[theItem.Id][x].OTP)
	}
	json.NewEncoder(w).Encode(theItem)

}

func createDebitTx(amt float64, item item.Item)[]transaction.BaseTransaction{
	var txs []transaction.BaseTransaction
	var Debit transaction.BaseTransaction
	var PrvKeys []*ecdsa.PrivateKey
	var total float64
	
	
	if amt >100{
		numTx := int(amt)/100
		fmt.Println("Num Tx: =", numTx)
		for x :=0; x < numTx; x+=1{
			Debit.ChainYear = common.ChainYear
			Debit.BlockNumber = common.BlockNumber
			Debit.Time = time.Now()
			Debit.Amount = common.USD2FG(float64(100))
			total += float64(100)
			PrvKey, _ := crypto.GenerateKey()
			PubKey := &PrvKey.PublicKey
			PrvKeys = append(PrvKeys, PrvKey)
			Debit.OTP = crypto.EncodePubKey(PubKey)
			Debit.HashBaseTx(Debit.OTP)
			txs = append(txs, Debit)
		}
		leftOver := amt -total
		if leftOver > 0{
			Debit.ChainYear = common.ChainYear
			Debit.BlockNumber = common.BlockNumber
			Debit.Time = time.Now()
			Debit.Amount = common.USD2FG(leftOver)
			PrvKey, _ := crypto.GenerateKey()
			PubKey := &PrvKey.PublicKey
			PrvKeys = append(PrvKeys, PrvKey)
			Debit.OTP = crypto.EncodePubKey(PubKey)
			Debit.HashBaseTx(Debit.OTP)
			txs = append(txs, Debit)
		}
		
	}else{
		numTx := int(amt)/10
		for x :=0; x < numTx; x+=1{
			Debit.ChainYear = common.ChainYear
			Debit.BlockNumber = common.BlockNumber
			Debit.Time = time.Now()
			Debit.Amount = common.USD2FG(10)
			total += float64(10)
			PrvKey, _ := crypto.GenerateKey()
			PubKey := &PrvKey.PublicKey
			PrvKeys = append(PrvKeys, PrvKey)
			Debit.OTP = crypto.EncodePubKey(PubKey)
			Debit.HashBaseTx(Debit.OTP)
			txs = append(txs, Debit)
		}
		leftOver := amt -total
		if leftOver > 0{
			Debit.ChainYear = common.ChainYear
			Debit.BlockNumber = common.BlockNumber
			Debit.Time = time.Now()
			Debit.Amount = common.USD2FG(leftOver)
			PrvKey, _ := crypto.GenerateKey()
			PubKey := &PrvKey.PublicKey
			PrvKeys = append(PrvKeys, PrvKey)
			Debit.OTP = crypto.EncodePubKey(PubKey)
			Debit.HashBaseTx(Debit.OTP)
			txs = append(txs, Debit)
		}
	}
	fmt.Println("Number of Txs:=", len(txs))
	
	common.Wallet.Items.Keys[item.Id] = PrvKeys
	
	return txs
}

func GetWallet(w http.ResponseWriter, r *http.Request){
	common.Wallet.FGValue = common.FGValue

	json.NewEncoder(w).Encode(common.Wallet)
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
			if !common.SendTransaction(tx){
				bn := common.BlockNumber + uint64(1)
				common.BlockFailed(bn)
			}
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

func NewNode() node.Node{
	var node node.Node
	node.PrvKey, _ = crypto.GenerateKey()
	node.PubKey = &node.PrvKey.PublicKey
	node.Id = crypto.GetAddress(node.PubKey) 
	node.PRKStr, node.PKStr  = crypto.Encode(node.PrvKey, node.PubKey)
	node.Port = ":"+ port
	node.Path = path
	//fmt.Println("Ip Address:", ipAddress)
	node.Tor = torPort
	node.Leader = false
	node.Writer =false
	return node
}


func newBlock(){
	common.CreateBlock(  )
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
	
	txDebitAmt.Add(txDebitAmt, Tx.Debit.Amount)

	txInt := Tx.CalcInterest()
	txFee := Tx.CalcFee()
	txCAmt.Add(txInt, txCreditAmt)
	txDAmt.Add(txFee, txDebitAmt)
	Tx.Change.Amount.Sub(txCAmt, txDAmt)
	txAcc.Add(txDAmt,Tx.Change.Amount)
	txAcc.Sub(txAcc,txCAmt)
	if (txAcc.Cmp(big.NewInt(0))==0){

		if bytes.Compare([]byte(Tx.TxHash), []byte(Tx.HashTx())) ==0 {
			if Tx.Payout == false{
				x:=0;
				for _,V :=range common.Wallet.Items.Keys{
					pubKey := crypto.EncodePubKey(&V[x].PublicKey)
					if pubKey == Tx.Debit.OTP{
						common.Wallet.Debits.Debit = map[string]transaction.BaseTransaction{Tx.Debit.TxHash: Tx.Debit}
						common.Wallet.FGs += common.Wei2FG(Tx.Debit.Amount)
					}
					x +=1
				}
				return Tx
				}
			}
	}
	return inValidTx
}

