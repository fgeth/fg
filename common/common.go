package common

import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "io/ioutil"
	 "math/big"
	 "os"
	 "path/filepath"
	 "strconv"
	 "sync"
	 "time"
	 "github.com/fgeth/fg/block"
	 "github.com/fgeth/fg/chain"
	 "github.com/fgeth/fg/crypto"
	 "github.com/fgeth/fg/item"
	 "github.com/fgeth/fg/node"
	 "github.com/fgeth/fg/transaction"
)

var (
	ChainYear			uint64							//Current Year
	BlockNumber			uint64							//Current Block Number
	ActiveNodes			[]string						//Array of known active Nodes Public Key as string
	PB					*block.Block						//Current Block This is the Last Know Verified Block
	Tx					[]transaction.Transaction		//Last Known Verified Block Transactions
	PBTx				[]transaction.Transaction		//Previous Block Transactions
	BTx					[]transaction.Transaction		//Used to Store Transactions for Pending Block
	PTx					[]crypto.Hash					//Array of Transaction Hashes for Pending Block
	Chain				chain.Chain						//Current Chain
	Chains				chain.Chains					//All Past Year Chains 
	FGValue				float64							//The Value of 1 FG
	Active				[]node.Node						//All Known Active Nodes Next Block
	Nodes				node.Nodes						//All known Nodes
	Writers				[]string						//Array of Current Block Nodes PublicKey as string Based on Block Hash includes Leader wich is the first node listed
	BTxHash				[]crypto.Hash					//Stores processed transaction debit hashes while Block or Leader Node
	PBTxHash			[]crypto.Hash					//Stores previous Block Transactions to account for Transactions sent to Block Leader until block is created & used to validate transactions are in block
	NumTx				int64							//Keeps track of number of Transactions resets at 1,000 Transactions and FGValue is bumped .01
	TTx					[]transaction.Transaction	    //Used to Transfer Transactions To Nodes One Block at a Time	
	Items				map[string]item.Item			//Index is Item Id
	MyNode				node.Node
	Mtx	 				sync.Mutex
	Path				string							//Path to Data dirctory
)

//Increments ChainYear by one
func IncChainYear(){
	Mtx.Lock()
	ChainYear += uint64(1)
	Mtx.Unlock()
	fmt.Println(ChainYear)
}



//Swap Pervious Block with now Current Block
func SwapBlocks(block *block.Block){
	Mtx.Lock()
	BlockNumber +=  uint64(1)
	PB = block
	Mtx.Unlock()
	fmt.Println(BlockNumber)
}

//Add Transaction to Last Known Verified Block Transactions
func SwapTransaction(){
	Mtx.Lock()
	PBTx = Tx
	Tx = BTx
	BTx = nil
	Mtx.Unlock()

}

//Trims Pending block transaction hashes since htere were over 1,000 trnasactions these need to be move to the next block
func trimPTx(){
	Mtx.Lock()
	if len(PTx) > 1000{
		var tmpTx  []crypto.Hash
		for x :=1000; x < len(PTx); x+=1{
			tmpTx = append(tmpTx, PTx[x])
		}
		
		PTx = tmpTx
		
	}
	Mtx.Unlock()
	time.Sleep(2 * time.Second) 
	CreateBlock()	
}


//Adds Transaction to Pending Block Transaction array
func AddBTX(tx transaction.Transaction){
	Mtx.Lock()
	BTx	 = append (BTx, tx)
	Mtx.Unlock()
}

//Swap out Active Nodes Array
func SwapActiveNodes(an []string){
	Mtx.Lock()
	ActiveNodes = an
	Mtx.Unlock()

}

func FG2USD(amount *big.Int) float64{
    fg := new(big.Int)
	fg.SetString("1000000000000000000", 10)

	f := new(big.Float).SetInt(amount)
	t := new(big.Float).SetInt(fg)
	f = f.Quo(f, t)

	fv, _:= f.Float64()
	
	usd :=   FGValue * fv
	return usd
	
	
}

func USD2FG(amount float64) *big.Int{
	
	bigval := new(big.Float)

	fgs := amount / FGValue

	bigval.SetFloat64(fgs)

	fv := new(big.Float)
	fv.SetString("1000000000000000000")

	fg :=new(big.Int)

	bigval.Mul(bigval, fv)
	bigval.Int(fg)

	return fg
	
	
}

func CreateBlock( ) block.Block{
var block  block.Block
if MyNode.Leader{
	blockNumber:= BlockNumber + uint64(1)
	NumTxs := uint64(len(PTx))
	var blockTx  []crypto.Hash
	TxFees :=big.NewInt(0)
	MultiBlock :=false
	if NumTxs > 1000{
		MultiBlock = true
		for x:=0; x < 1000; x +=1{
			blockTx = append (blockTx, PTx[x])
			percentage := big.NewInt(100)
			for y:=0; y< len(BTx[x].Credit); y +=1{
				txFee := new(big.Int).Div(BTx[x].Credit[y].Amount, percentage)
				TxFees = TxFees.Add(TxFees,txFee)
			}
		}
		
	}else{
		blockTx = PTx
	}
	BlockReward :=big.NewInt(0)
	bn := int64(len(blockTx))
	NumTx += bn
	if bn <1000{
		n:= big.NewInt(0)
		n.SetString("10000000000000000", 10)
		t:= big.NewInt(bn)
		BlockReward =new(big.Int).Mul(t, n)
	}else{
		BlockReward= big.NewInt(0)
		BlockReward.SetString("10000000000000000000", 10)
	}
		fmt.Println(BlockReward)
	
	if NumTx > 1000 {
		if FGValue <1000{
			FGValue +=.01 
			
		}
		if (FGValue >=1000) &&(FGValue < 10000){
			FGValue +=.001
			s:=big.NewInt(2)
			BlockReward = new(big.Int).Div(BlockReward, s)

		}
		if FGValue >=10000{
			if FGValue <100000{
				FGValue +=.0001
				b:=big.NewInt(10)
				BlockReward = new(big.Int).Div(BlockReward, b)
			}else{
				b:=big.NewInt(50)
				BlockReward = new(big.Int).Div(BlockReward, b)
			
			}
		}
		
		NumTx = NumTx - 1000
		
	}
	
	NodeTx := PayOutNodes(TxFees, blockNumber)	
	for x:=0; x < len(NodeTx); x +=1{
			blockTx = append (blockTx, NodeTx[x].TxHash)
			AddBTX(NodeTx[x])
			
			
		}
	WritersTx := PayOutWriters(BlockReward, blockNumber)	
	for x:=0; x < len(WritersTx); x +=1{
			blockTx = append (blockTx, WritersTx[x].TxHash)
			AddBTX(WritersTx[x])
		}
		
		t := time.Now()
		if uint64(t.Year())> ChainYear{
			IncChainYear()
			block.ChainYear = ChainYear
		}
		block.ChainYear = ChainYear
		block.BlockNumber = blockNumber
		block.FGValue = FGValue
		block.Txs = blockTx
		block.NumTxs = NumTxs
		block.Nodes = ActiveNodes
		block.PBHash = PB.BlockHash
		block.NodePayout = NodeTx[0].Debit[0].Amount
		block.WriterPayout = WritersTx[0].Debit[0].Amount
		block.BlockHash = block.HashBlock()
		if MultiBlock{
			block.Writers = PB.Writers
			go trimPTx()
		}else{
			nodeVals := ElectNodes(block)
			for x:=0; x < len(nodeVals); x +=1{
				block.Writers = append(block.Writers, block.Nodes[nodeVals[x]])
			}
		}
			
		SwapBlocks(&block)
		block.SaveBlock(MyNode.Path)
	 
	}
	return block
}


//TODO Recreate Block Failed
func BlockFailed(blockNumber uint64){

}


func ElectNodes(block block.Block) []uint64{

	numTx := (PB.NumTxs/500)+1
	var blockNode []uint64
	var numNodes uint64
	numNodes = uint64(len(block.Nodes))
	if numNodes < 7{
		numTx = numNodes
	}else{
		if numTx < 7 {
			numTx =7
		}else{
			if numTx > numNodes{
				numTx = numNodes
			}
		}
	}
	if bytes.Compare(block.BlockHash, block.HashBlock()) ==0 {
	    uintA, uintB, uintC, uintD :=crypto.HashToUint64(block.BlockHash)
		
		tmp := uint64(0)
		bn := uint64(0)
		for x :=uint64(0); x < numTx; x +=1{
				switch x%4{
					case  1:
						tmp = (uintA % numNodes)+x
						
							
					case 2:
						tmp = (uintB % numNodes)+x
					
					case 3:
						tmp = (uintC % numNodes)+x
					
					case 4:
						tmp = (uintD % numNodes)+x
		
				}	
					if tmp > numNodes{
						bn = tmp - numNodes
					}else{
						bn = tmp
					}
					blockNode  = append(blockNode, (bn))
				
			}
			
			
			
	}
	return blockNode
}

func VerifyBlock(block *block.Block) bool{
	numNodes := 0
	for x:=0; x < len(block.Signed); x +=1{
		if block.VerifyBlock(PB){
				if (len(block.Txs) ==1000) && (CompareWriters(block.Writers,PB.Writers)){
					if block.FGValue - PB.FGValue <=.01{
						numNodes +=1
					}
				}else{
					
					if CompareWriters(block.Writers, block.GetWriters(ElectNodes(*block))){
						numNodes +=1
					}
				}
			
		}
	}
		//Verify that Block Leader and a Block Node signed the Block
	if numNodes > 1{
		//TODO Store Block to file, Replace Previous Block wtih This Block, 
		if bytes.Compare(block.BlockHash, block.HashBlock()) ==0{	
			
			if bytes.Compare(block.PBHash, PB.BlockHash) ==0{
				if CheckBlock(block){
					CB := ImportBlock(block)		// current saved block need to add the other block signatures
					for x:=0; x < len(CB.Signed); x +=1{
						if CB.VerifyBlock(PB){
							block.Signed = append(block.Signed, CB.Signed[x])
						}
					}
				}
				block.SaveBlock(MyNode.Path)
				if BlockNumber < block.BlockNumber{
					SwapBlocks(block)
				}
				return true
			}
		}
		
	}
	return false

}

func CheckBlock(block *block.Block) bool{
	dirname, err := os.UserHomeDir()

    if err != nil {
        fmt.Println( err )
    }
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "fg", "chain", strconv.FormatUint(block.ChainYear, 10))

	fileName := filepath.Join(path, strconv.FormatUint(block.BlockNumber, 10))
	myfile, e := os.Stat(fileName)
	if e != nil{
	  fmt.Println( e )
	  return false
	}
	fmt.Println( myfile )
	return true

}

func ImportBlock(CB *block.Block) block.Block{
	return block.ImportBlock(CB.ChainYear, CB.BlockNumber)

}

func CompareWriters(writers []string, theWriters []string) bool{
	NotEqual := false
	for x :=0; x < len(theWriters); x +=1{
		if writers[x] != theWriters[x]{
			NotEqual = true
		}
			
	}
	return NotEqual
}

func PayOutNodes(TxFees *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(int64(len(ActiveNodes)))
payOut := TxFees.Div(TxFees, nodes)

var Txd []transaction.Transaction
	for k,x := range ActiveNodes{
		Txd =append(Tx, CreatePayoutTransaction(payOut, x, blockNumber))
		fmt.Println("Paying node number:", k)
	}
	
	return Txd

	
}

func PayOutWriters(blockReward *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(int64(len(Writers)))
payOut := blockReward.Div(blockReward, nodes)

var Txd []transaction.Transaction
	for x:=0; x<len(Writers); x +=1{
		Txd=append(Tx, CreatePayoutTransaction(payOut, Writers[x], blockNumber))
	}
	
	return Txd

	
}

func CreatePayoutTransaction(amt *big.Int, pubKey string, blockNumber uint64) transaction.Transaction{


	var Debit transaction.BaseTransaction
	var Credit transaction.BaseTransaction
	var Tx transaction.Transaction
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
	Tx.OTP = MyNode.Id
	Tx.TxHash = Tx.HashTx()
	Tx.R, Tx.S = crypto.Sign(Tx.TxHash, MyNode.PrvKey)
	Tx.Payout = true
	return Tx
	
}
func SellItem(item item.Item) {
	//prvKey := crypto.GenerateRSAKey()
	//MyNode.Comms.RsaPrvKeys[prvKey.PublicKey] = prvKey
	//item.Seller = prvKey.PublicKey
	//MyNode.Items.Item = append(MyNode.Items.Item, item)

	
}

//TODO Fix To Where this imports the blocks
func ImportBlocks() {
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "fg", "chain", strconv.FormatUint(ChainYear, 10))

	fileName := filepath.Join(path, strconv.FormatUint(BlockNumber, 10))
	file, _ := ioutil.ReadFile(fileName)
	var block block.Block
 	_ = json.Unmarshal([]byte(file), &block)

}
//TODO Sign Genesis block
func SignGenesisBlocks(){

}

//TODO GetBlocks
func GetBlocks(){

}
//TODO Get Txs
func GetTxs(){

}

func VaildTransaction(Tx transaction.Transaction) bool{
txFee := Tx.CalcFee()
 txInterest := big.NewInt(0)
if FGValue >=100{
	txInterest = Tx.CalcInterest()
}

if Tx.OTP !=""{
		if Tx.VerifySig(){
			TC := Tx.Credits()
			TD := Tx.Debits()
			if TC.Add(TC,txInterest) == txFee.Add(txFee, TD.Add(TD, Tx.Change.Amount)){
			for x:=0; x< len(Tx.Credit); x+=1{
				cTx := transaction.ImportBaseTx(Tx.Credit[x].TxHash)
				if cTx.OTP ==""{
					cTx.OTP=Tx.Credit[x].OTP
					cTx.SaveTx()
					
				}else{
					return false		//Double Spend
				}
			}
			for x:=0; x< len(Tx.Debit); x+=1{
				Tx.Debit[x].SaveTx()

			}
				Tx.Change.SaveTx()
			}else{
					return false
			}
		}else{
			return false
		}
}else{
	if Tx.Payout == true{
		if Tx.VerifySig(){
			block := block.ImportBlock(Tx.Credit[0].ChainYear , Tx.Credit[0].BlockNumber)
			for x:=0; x<len(block.Writers); x +=1{
				if block.Writers[x] == Tx.OTP{
					if Tx.Debit[0].Amount == block.NodePayout{
							return true
						}
					if Tx.Debit[0].Amount == block.WriterPayout{
						return true
					}
				}
			}
		}
	
	}
	return false
}
return false
}



//TODO fix where this imports the transactions
func ImportTxs() {
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
	var txHash crypto.Hash
	uintA, uintB, uintC, uintD := crypto.HashToUint64(txHash)
	h1 := strconv.FormatUint(uintA, 10)
	h2 := strconv.FormatUint(uintB, 10)
	h3 := strconv.FormatUint(uintC, 10)
	h4 := strconv.FormatUint(uintD, 10)
	theHash := h1 + h2 +h3 +h4
	path :=filepath.Join(dirname, "fg", "tx", theHash )
	file, _ := ioutil.ReadFile(path)
	var tx transaction.Transaction
	_ = json.Unmarshal([]byte(file), &tx)
}




func AllItemsInDir() {
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }

	dir :=filepath.Join(dirname, "fg", "items")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		  Items[f.Name()] = item.ImportItem(f.Name())
	}
   
	  
}

func genRsa() {
 prvKey := crypto.GenerateRSAKey()
 pubKey := prvKey.PublicKey
 secret := crypto.RSAEncrypt("Secret", pubKey)
 fmt.Println("Encrypted Message =", secret)
 clearText := crypto.RSADecrypt(secret, prvKey)
 fmt.Println("Message =", clearText)
 publicKey := crypto.EncodeRSAPubKey(&pubKey)
 fmt.Println("Publick Key =", publicKey)
 pKey := crypto.DecodeRSAPubKey(publicKey)
 secret2 := crypto.RSAEncrypt("Secret2", pKey)
 fmt.Println("Encrypted Message =", secret2)
 clearText2 := crypto.RSADecrypt(secret2, prvKey)
 fmt.Println("Message =", clearText2)
 pk := prvKey
 err:= crypto.StoreRSAKey( pk ,"Pass", "Key1")
 fmt.Println(err)
 pvKey, pbKey, err :=crypto.GetRSAKey("Key1", "Pass") //rsa.PrivateKey,rsa.PublicKey, error
 secret3 := crypto.RSAEncrypt("Secret3", pbKey)
 fmt.Println("Encrypted Message =", secret3)
 clearText3 := crypto.RSADecrypt(secret3, pvKey)
 fmt.Println("Message =", clearText3)


}