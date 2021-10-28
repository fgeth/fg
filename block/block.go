package block

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	 "os"
	"path/filepath"
	"strconv"
	"github.com/fgeth/fg/crypto"


)

type Block struct {
	ChainYear			uint64						//Year this block was created
	BlockNumber			uint64						//Block Number
	FGValue				float64						//The Value of 1 FG
	Txs			  		[]string					//Array of Transaction Hashes
	NumTxs				uint64						//Number of Transactions Submited For this Block
	Nodes				[]string					//This is current list of Nodes that responded to the last Block.  This array is what is used to determine Block Nodes. string is the nodes public key as a string
	PBHash				crypto.Hash					//Hash of previous Block
	BlockHash			crypto.Hash					//Hash of this Block which includes previous Blocks Hash & list of Nodes
	Writers				[]string					//Array of the Next Block Nodes PublicKey as string Based on Block Hash includes Leader
	Signed				[]SignedBlock				//Signature of Block Nodes
	BlockFailed			bool						//Used if Block fails to be created something happens to all Block Nodes
	NodePayout			*big.Int						//Amount paid out to each Acitve Nodes when Block is created
	WriterPayout		*big.Int						//Amount paid out to each Blcok Writer includes Block Leader when Blcok is Created
	BlockReward			*big.Int						//Amount of FG in Block Reward paid to Block Writers and Leader
	
}

type SignedBlock struct {
	R			*big.Int
	S			*big.Int
	NodeId		string

}

type MinBlock struct{
	ChainYear			uint64						//Year this block was created
	BlockNumber			uint64						//Block Number
}
func (block *Block) GetUnsignedBlock() Block{

	var unsigned Block
	unsigned.ChainYear = block.ChainYear
	unsigned.BlockNumber = block.BlockNumber
	unsigned.FGValue = block.FGValue
	unsigned.Txs = block.Txs
	unsigned.NumTxs = block.NumTxs
	unsigned.Nodes = block.Nodes
	unsigned.NodePayout = block.NodePayout
	unsigned.WriterPayout = block.WriterPayout
	unsigned.BlockFailed = block.BlockFailed
	unsigned.PBHash = block.PBHash
	
 return unsigned
}

func (block *Block) HashBlock() crypto.Hash{
	kh :=crypto.NewKeccakState()
	unsignedBlock := block.GetUnsignedBlock()
	json , _:= json.Marshal(unsignedBlock)
	
	blockHash :=crypto.HashData(kh, []byte(json))
	return blockHash
}


func (block *Block) SaveBlock(dirname string){
	//dirname, err := os.UserHomeDir()
   // if err != nil {
     //   fmt.Println( err )
    //}
 
	path :=filepath.Join(dirname, "chain", strconv.FormatUint(block.ChainYear, 10))
	 
	_, err := os.Stat(path)
	
    if os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		fmt.Println(err)
        err1 := os.Mkdir(filepath.Join(dirname, "chain"), 0755)
		fmt.Println(err1)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
  
	fileName := filepath.Join(path,strconv.FormatUint(block.BlockNumber, 10))
	file, _ := json.MarshalIndent(block, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}

func ImportBlock(chainYear uint64, blockNumber uint64, dirname string) Block{
	//dirname, err := os.UserHomeDir()
	var block Block
    //if err != nil {
     //   fmt.Println( err )
    //}
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "chain", strconv.FormatUint(chainYear, 10))

	fileName := filepath.Join(path, strconv.FormatUint(blockNumber, 10))
	myfile, e := os.Stat(fileName)
	if e != nil{
	  fmt.Println( e )
	}else{
		file, _ := ioutil.ReadFile(myfile.Name())
		
		_ = json.Unmarshal([]byte(file), &block)
		
		
	}
	return block
}

func (block *Block) GetWriters(nodes []uint64) []string{
	var writers []string
	for x:=0; x < len(nodes); x +=1{
		writers = append(writers, block.Nodes[nodes[x]])
	}
	return writers
}
func (block *Block) VerifyBlock(PB *Block) bool{
	return block.VerifyWriters(PB)

}


func (block *Block) VerifyWriters(PB *Block) bool{
	Signed :=0
	NumWriters := len(PB.Writers)
	for x:=0; x < NumWriters; x +=1{
		for w:=0; w < NumWriters; w +=1{
		if block.Signed[x].NodeId == PB.Writers[w]{
			if block.VerifySig(w){
				Signed +=1
			}
			
		}

	}
	if Signed >= NumWriters /2{
		return true
	}
	}
return false
}

func (block *Block) VerifySig(index int) bool{

	blockHash := block.HashBlock()
	publicKey := crypto.DecodePubKey(block.Signed[index].NodeId)
	if bytes.Compare(block.BlockHash, blockHash) ==0 {
		return crypto.Verify(block.BlockHash, block.Signed[index].R, block.Signed[index].S, publicKey)
	}else{
		return false
	}
}



