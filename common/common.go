package common

import (
	 "bytes"
	 "crypto/ecdsa"
	 "encoding/binary"
	 "fmt"
	 //"io/ioutil"
	 //"crypto/ecdsa"
	 "math/big"
	 
	 "net/http"
	//"net/http/cookiejar"

	 //"strconv"
	 "sync"
	 "time"
	 "github.com/fgeth/fg/block"
	 "github.com/fgeth/fg/chain"
	 "github.com/fgeth/fg/crypto"
	 "github.com/fgeth/fg/item"
	 "github.com/fgeth/fg/node"
	 "github.com/fgeth/fg/ring"
	 "github.com/fgeth/fg/transaction"
	 "github.com/fgeth/fg/wallet"
	 
)

var (
	ChainYear			uint64							//Current Year
	Ring				ring.Ring						//Block Ring
	BlockNumber			uint64							//Current Block Number
	ActiveNodes			[]string						//Array of known active Nodes Public Key as string
	PB					*block.Block						//Current Block This is the Last Know Verified Block
	Tx					[]transaction.Transaction		//Last Known Verified Block Transactions
	PBTx				[]transaction.Transaction		//Previous Block Transactions
	BTx					[]transaction.Transaction		//Used to Store Transactions for Pending Block
	PTx					[]string						//Array of Transaction Hashes for Pending Block
	Chain				chain.Chain						//Current Chain
	Chains				chain.Chains					//All Past Year Chains 
	FGValue				float64							//The Value of 1 FG
	Active				[]node.Node						//All Known Active Nodes Next Block
	TheNodes			 Nodes							//All known Nodes
	Writers				[]string						//Array of Current Block Nodes PublicKey as string Based on Block Hash includes Leader wich is the first node listed
	BTxHash				[]string						//Stores processed transaction debit hashes while Block or Leader Node
	PBTxHash			[]string						//Stores previous Block Transactions to account for Transactions sent to Block Leader until block is created & used to validate transactions are in block
	NumTx				int64							//Keeps track of number of Transactions resets at 1,000 Transactions and FGValue is bumped .01
	TTx					[]transaction.Transaction	    //Used to Transfer Transactions To Nodes One Block at a Time	
	Items				map[string]item.Item			//Index is Item Id
	MyNode				node.Node
	Mtx	 				sync.Mutex
	Path				string							//Path to Data dirctory
	Trusted				[]*ecdsa.PublicKey				//PublicKey of Fgeth Servers
	Wallet				wallet.Wallet
	OC					OnionClient
	Cookies         	[]*http.Cookie
)

type Nodes struct {
	Node 			map[string]node.Node  				//All active and inactive Nodes.  Easy to get PublicKey using string of public key as map index 

}


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
 func Byte2Uint64 (byteArray []byte) uint64{
	var ret uint64
	buf := bytes.NewBuffer(byteArray)
    binary.Read(buf, binary.BigEndian, &ret)
	return ret
 }
//Trims Pending block transaction hashes since htere were over 1,000 trnasactions these need to be move to the next block
func trimPTx(){
	Mtx.Lock()
	if len(PTx) > 1000{
		var tmpTx  []string
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

func Wei2FG(amount *big.Int) float64{
    fg := new(big.Int)
	fg.SetString("1000000000000000000", 10)

	f := new(big.Float).SetInt(amount)
	t := new(big.Float).SetInt(fg)
	f = f.Quo(f, t)

	fv, _:= f.Float64()
	
	return fv
	
	
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

