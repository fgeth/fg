package main

import(
	"bytes"
	//"context"
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
	//"strconv"
	"sync"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/node"
	"github.com/fgeth/fg/ring"
	"github.com/fgeth/fg/transaction"
	//"github.com/fgeth/bine/tor"
	//"github.com/fgeth/fasthttp"


)
var (

	wg 		 sync.WaitGroup
	path	  string
	port 	  string
	ipAddress string 		//Port Tor is running on
	auth	  string
	Gen		  *bool
	
	
)

func init() {
			flag.StringVar(&port, "port", "42069", "Default Port")
			flag.StringVar(&path, "path", "/var/fg", "Data Directory")
			flag.StringVar(&ipAddress, "ip", "127.0.0.1", "Default Port")
			flag.StringVar(&auth, "auth", "P@Ssw0Rd1", "Default Password")
	Gen = 	flag.Bool("gen", false, "Continue with existing chain")
	
}

func main(){	
	flag.Parse()
	common.Auth = auth
	var tmpNode node.Node
	tmpNode.Id =uint64(0)
	tmpNode.Ip =""
	tmpNode, err := node.ImportNode(path)
	
	if err !=nil{
		common.MyNode = NewNode()

	}else{
		fmt.Println("Found Node", tmpNode)
		common.MyNode = tmpNode
		
		if common.MyNode.PRKStr !=""{
			common.MyNode.PrvKey = crypto.DecodePrv(common.MyNode.PRKStr)
			fmt.Println("Private Key Added from saved Node")
			if common.MyNode.PKStr !=""{
				common.MyNode.PubKey = crypto.DecodePubKey(common.MyNode.PKStr)
				fmt.Println("Public Key Added from saved Node")
			}
		}else{
			fmt.Println("Error Getting Private Key")
			common.MyNode.PrvKey, common.MyNode.PubKey, err =crypto.GetKey(common.MyNode.Address, auth)
			if err !=nil{
				fmt.Println("Error Recovering Keys Generating new Keys")
				common.MyNode.PrvKey, _ = crypto.GenerateKey()
				common.MyNode.PubKey = &common.MyNode.PrvKey.PublicKey
				common.MyNode.PRKStr, common.MyNode.PKStr  = crypto.Encode(common.MyNode.PrvKey, common.MyNode.PubKey)
				common.MyNode.Address, _ = crypto.StoreKey( common.MyNode.PrvKey, auth, path)
			}
		}
	}
		fmt.Println("Node Id is :" , common.MyNode.Id)
		fmt.Println("Node Ip is :" , common.MyNode.Ip)
		fmt.Println("Node Address is :" , common.MyNode.Address)
		fmt.Println("Node Path is :" , common.MyNode.Path)
	common.MyNode.SaveNode(common.MyNode.Path)
	common.TheNodes.Node = map[string]node.Node{common.MyNode.PKStr: common.MyNode}
	tmpRing, err := ring.ImportRing(common.MyNode.Path)
	if err != nil{
		common.Ring = NewRing()
		finger := ring.FingerTable{Id : uint64(0), Node: NodeOne()}
		common.Ring.Table = append(common.Ring.Table, finger)
		finger = ring.FingerTable{Id : uint64(1), Node : NodeTwo()}
		common.Ring.Table = append(common.Ring.Table, finger)

		fmt.Println("Node One Pub Key From Ring ", common.Ring.Table[0].Node.PKStr)
		fmt.Println("Node Two Pub Key From Ring ", common.Ring.Table[1].Node.PKStr)
		common.Ring.SaveRing(common.MyNode.Path)
	}else{
		common.Ring = tmpRing
	}
	
	
	if *Gen{
		fmt.Println("Genesis Block")
		common.FGValue = .01
		common.CreateGenBlocks()
		
		//common.ImportTx()
		//common.SignGenesisBlocks() 
	}else{
		//RegisterNode()
	
		if common.MyNode.Id == 0{
		//	common.Ring.FindPeer()
			// RegisterNode()
		}else{
			fmt.Println("Node Id: ", common.MyNode.Id)
		}
		Trusted()
		if port !="42069"{
			common.MyNode.Port =":"+ port
			}
		if ipAddress !="127.0.0.1"{
			common.MyNode.Ip = ipAddress
		}
	
		directory()
		wallet()
		
		
		test()
		//common.ImportBlocks()
		//common.ImportTxs()
		//common.GetBlocks()
		//common.GetTxs()
	}
	wg.Add(1)
	go server()
	go postTest()
	//torServer()
	
	
	//TorService()
	//time.Sleep(time.Second * 120)
	
	go fg()
	go CloseHandler()
	wg.Wait()
	
}

func NewRing() ring.Ring{
	var finger []ring.FingerTable
	var rnodes []node.RNode
	ring := ring.Ring {uint64(0), finger, rnodes }
	return ring
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
func postTest(){
	time.Sleep(5 * time.Second)
	
	common.ImportBlocks(uint64(0))
	block := common.GetBlock(0)
	fmt.Println("Block Received Over network")
	fmt.Println("BlockNumber:", block.BlockNumber)
	fmt.Println("ChainYear:", block.ChainYear)
	fmt.Println("FGValue:", block.FGValue)
	fmt.Println("Txs:", block.Txs)
	

}
func test(){

	common.CreateGenBlocks()
	
	//fmt.Println("Blocks :", common.Chain.Blocks[0])
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
	r.HandleFunc("/getBlock", sendBlock).Methods("POST")
	//r.HandleFunc("/getNodes", sendNodes).Methods("GET")
	//r.HandleFunc("/getTxs", sendTxs).Methods("GET")
	r.HandleFunc("/getWallet", GetWallet).Methods("GET")
	r.HandleFunc("/getPeer", GetPeer).Methods("GET")
	r.HandleFunc("/sendTx", sendNewTransaction).Methods("POST")
	//r.HandleFunc("/Tx", CreateNewTransaction).Methods("POST")
	r.HandleFunc("/block", createNewBlock).Methods("POST")
	r.HandleFunc("/addItem", createNewItem).Methods("POST")
	r.HandleFunc("/newNode", newNode).Methods("POST")
	//r.HandleFunc("/blockTxs", processTxs).Methods("POST")
	
	
	
	staticFileDirectory := http.Dir(path)
	
	staticFileHandler := http.StripPrefix("/store/", http.FileServer(staticFileDirectory))
	
	r.PathPrefix("/store/").Handler(staticFileHandler).Methods("GET")
	
	fmt.Println("Listening on port :", common.MyNode.Port)
	server := &http.Server{
    Addr:    common.MyNode.Port,
    Handler: r,
		}
if err := server.ListenAndServe(); err != nil {
	    	log.Fatal(err)
	   }
	   
}


//Function to check that new blocks are being made and if not start the process
func fg(){
	wg.Add(1)
	defer wg.Done()
	for {
		CheckBlockNumber := common.BlockNumber	
		//time.Sleep(time.Second * 60)
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

func sendBlock(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Request from Network:", reqBody)
	var minBlock block.MinBlock
    json.Unmarshal(reqBody, &minBlock)
	//var chainYear uint64
    //chainYear = common.Byte2Uint64(KV.Args[1].Value)
	//var blockNum uint64
   // blockNum = common.Byte2Uint64(KV.Args[0].Value)
	fmt.Println("Block from Network:", minBlock)
	Block := block.ImportBlock(minBlock.ChainYear, minBlock.BlockNumber, common.MyNode.Path)
	fmt.Println("Block from File:", Block)
	json.NewEncoder(w).Encode(Block)
}

func newNode(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
    var theNode node.SNode
    json.Unmarshal(reqBody, &theNode)
	var rnode node.RNode
	rnode.Id,_,_,_ = crypto.B32HashToUint64([]byte(crypto.HashTx([]byte(theNode.PKStr))))
	rnode.PKStr = theNode.PKStr
	common.Ring.RotateKeys(rnode)
	common.Ring.RotateFingerTable(theNode, common.MyNode.Id)
	json.NewEncoder(w).Encode(theNode)
	
}

func RegisterNode( ) {
	nodeJson, _ := json.Marshal(common.MyNode)
	for x:=0; x < len(common.Ring.Table); x +=1{
		url := "http://"+common.Ring.Table[x].Node.Ip+common.Ring.Table[x].Node.Port+"/newNode"
		fmt.Println("Connecting to Ring at :", url)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(nodeJson))

		if err != nil {
			fmt.Println("Could not make POST request to ring trying port 80 if 42069 was banned by governments or ISPs")
			thePort := common.Ring.Table[x].Node.Port
			if ( thePort == ":42069"){
				url := "http://"+common.Ring.Table[x].Node.Ip+":80/newNode"
				fmt.Println("Connecting to Ring at :", url)
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(nodeJson))
				if err != nil {
					fmt.Println("Could not make POST request to ring")
				}else{
					body, err := ioutil.ReadAll(resp.Body)

					var result node.Node
					err = json.Unmarshal([]byte(body), &result)
					if err != nil {
						fmt.Println("Error unmarshaling data from request.")
					}else{
						common.MyNode.Id = result.Id
						common.Ring.Id = result.Id
						fmt.Println(" Registered Node Id: ", common.MyNode.Id)
						break
					}
				}
			}
		}else{

			body, err := ioutil.ReadAll(resp.Body)

			var result node.Node
			err = json.Unmarshal([]byte(body), &result)
			if err != nil {
				fmt.Println("Error unmarshaling data from request.")
			}else{
				common.MyNode.Id = result.Id
				common.Ring.Id = result.Id
				fmt.Println(" Registered Node Id: ", common.MyNode.Id)
				break
			}
		}
	}

	
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
	//path := filepath.Join(common.MyNode.Path, "Keys", theItem.Id) 
	crypto.StoreRSAKey(prvKey, "Password", theItem.Id, common.MyNode.Path)
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

func GetPeer(w http.ResponseWriter, r *http.Request){
	

	json.NewEncoder(w).Encode(common.MyNode)
}
//TODO Fix this
func register() bool{
	haveNode := false
	ImportActiveNodes()
	if len(common.ActiveNodes) >0{
		haveNode = true	
	}else{
		common.MyNode.GetNodes()
		if len(common.ActiveNodes) >0{
			haveNode = true	
		}
	}	
	if haveNode {
	for x:=0; x<len(common.ActiveNodes); x+=1{
			
			common.MyNode.RegisterNode(common.ActiveNodes[x])
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
	node.PRKStr, node.PKStr  = crypto.Encode(node.PrvKey, node.PubKey)
	node.Address, _ = crypto.StoreKey( node.PrvKey, auth, path)
	fmt.Println("Node Prvt Key File : ", node.Address)
	node.Port = ":"+port
	node.Path = path
	fmt.Println("Ip Address:", ipAddress)
	node.Ip = ipAddress
	node.Leader = false
	node.Writer =false
	return node
}


func NodeTwo() node.SNode{
	var node node.SNode
	node.Port = ":42069"
	node.Ip = "node2.fgeth.com"
	PrvKey, _ := crypto.GenerateKey()
	node.Address, _ =crypto.StoreKey(PrvKey, auth, path)
	node.PKStr  = crypto.EncodePubKey(&PrvKey.PublicKey)
	node.SaveNodeTwo("/var/fg")
	return node
}

func NodeOne() node.SNode{
	var node node.SNode
	node.Port = ":42069"
	node.Ip = "node1.fgeth.com"
	PrvKey, _ := crypto.GenerateKey()
	node.Address, _ =crypto.StoreKey ( PrvKey, auth, path)
	node.PKStr  = crypto.EncodePubKey(&PrvKey.PublicKey)
	node.SaveNodeOne("/var/fg")
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

