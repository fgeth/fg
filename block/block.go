package block

import (
	"time"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	 "os"
	"path/filepath"
	 "github.com/fgeth/fg/crypto"
	 "github.com/fgeth/fg/transaction"
)

type Block struct {
	ChainYear			uint64						//Year this block was created
	BlockNumber			uint64						//Block Number
	FGValue				float64						//The Value of 1 FG
	Txs			  		[]crypto.Hash				//Array of Transaction Hashes
	NumTxs				uint32						//Number of Transactions Submited For this Block
	Nodes				[]string					//This is current list of Nodes that responded to the last Block.  This array is what is used to determine Block Nodes. string is the nodes public key as a string
	PBHash				crypto.Hash					//Hash of previous Block
	BlockHash			crypto.Hash					//Hash of this Block which includes previous Blocks Hash & list of Nodes
	Writers				[]string					//Array of the Next Block Nodes PublicKey as string Based on Block Hash includes Leader
	Signed				[]SignedBlock				//Signature of Block Nodes
	BlockFailed			bool						//Used if Block fails to be created something happens to all Block Nodes
	NodePayout			big.Int						//Amount paid out to each Acitve Nodes when Block is created
	WriterPayout		big.Int						//Amount paid out to each Blcok Writer includes Block Leader when Blcok is Created
	
}

type SignedBlock struct {
	R			big.Int
	S			big.Int
	NodeId		string

}


func trimPTx(){
	if len(PTx) > 1000{
		tmpTx := []Transaction
		for x :=1000; x < len(PTx); x+=1{
			tmpTx = append(tmpTx, PTx[x])
		}
		PTx = tmpTx
	}
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
	unsigned.WriterPayout = block.writerPayout
	unsigned.BlockFailed = block.BlockFailed
	unsigned.PBHash = block.PBHash
	
 return unsigned
}

func (block *Block) HashBlock() crypto.Hash{
	kh crypto.NewKeccakState()
	unsignedBlock := block.GetUnsignedBlock()
	json := json.Marshal(unsignedBlock)
	
	blockHash :=crypto.HashData(kh, []byte(json))
	return blockHash
}

func(block *Block) ElectNodes() []uint64{
kh crypto.NewKeccakState()
	numTx := (PB.NumTxs/500)+1
	blockNode := []uint64
	numNodes := len(block.Nodes)
	if numNodes < 7{
		numTx := numNodes
	}else{
		if numTx < 7 {
			numTx =7
		}else{
			if numTx > numNodes{
				numTx = numNodes
			}
		}
	}
	blockHash == block.HashBlock(){
	    uintA, uintB, uintC, uintD :=crypto.HashToUint64(blockHash)
		
	
		for x :=0; x < numTx; x +=1{
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
			
			return blockNode
			
}
}


func (block *Block) VerifyBlock(){
	numNodes := 0
	for x:=0; x < len(block.Signed); x +=1{
		if block.VerifyWriters(){
				if len(block.Txs) ==1000 & (block.Writers == PB.Writers){
					if block.FGValue - PB.FGValue <=.1{
						numNodes +=1
					}
				}else{
					
					if block.Writers ==block.ElectNodes(){
						numNodes +=1
					}
				}
			
		}
	}
		//Verify that Block Leader and a Block Node signed the Block
	if numNodes > 1{
		//TODO Store Block to file, Replace Previous Block wtih This Block, 
		block.SaveBlock()
		
	}

}
func (block *Block) SaveBlock(){
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
 
	path :=filepath.Join(dirname, "fg", "chain", strconv.FormatUint(block.ChainYear, 10))
	 
	folderInfo, err := os.Stat(path)
	if folderInfo.Name() !="" {
			fmt.Println("")
	}
    if os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(dirname, "fg"), 0755)
		fmt.Println(err)
        err1 := os.Mkdir(filepath.Join(dirname, "fg", "chain"), 0755)
		fmt.Println(err1)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
  
	fileName := filepath.Join(path,strconv.FormatUint(block.BlockNumber, 10))
	file, _ := json.MarshalIndent(block, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}

func ImportBlock(chainYear uint64, blockNumber uint64) Block{
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "fg", "chain", strconv.FormatUint(chainYear, 10))

	fileName := filepath.Join(path, strconv.FormatUint(blockNumber, 10))
	file, _ := ioutil.ReadFile(fileName)
	var block Block
 	_ = json.Unmarshal([]byte(file), &block)
	
	return block
}


func(block *Block) VerifyWriters( ) bool{
	Signed :=0
	NumWriters := len(block.Writers)
	for x:=0; x < ; x +=1{
		for w:=0; w < NumWriters; w +=1{
		if block.Signed[x].NodeId == block.Writers[w]{
			if block.VerifySig(w){
				Signed +=1
			}
			
		}

	}
	if Signed => NumWriters /2{
		return TRUE
	}
	return False
}
func(block *Block) VerifySig(x uint32) bool{
	kh crypto.NewKeccakState()
	blockHash := block.HashBlock()
	publicKey := DecodePubKey(block.Signed[x].NodeId)
	if block.BlockHash == blockHash{
		return crypto.verify(block.Hash, R, S, publicKey)
	}else{
		return FALSE
	}
}



