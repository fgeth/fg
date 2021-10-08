package block

import (
	"time"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	 "os"
	 "github.com/fgeth/fg/common"
)

type Block struct {
	ChainId				uint						//Year of Block
	BlockNumber			uint64						//Block Number
	Txs					[]Transaction	   			//Array of Transactions in utc time order
	Nodes				[]uintptr					//This is current list of Nodes that responded to the last Block.  This array is what is used to get BlockNodes. 
	PBHash				common.Hash					//Hash of previous Block
	Writers				[]uintptr					//array of the Block Nodes uintptr  Based on Block Hash includes Leader
	BlockHash			common.Hash					//Hash of this Block which includes previous Blocks Hash & list of Nodes
	ChainHash			common.Hash					//Hash of all the Block hashes up to this point for this chain includes this BlockHash
	Signer				[]SignedBlock				//Signature of Block Nodes
	
}


type SignedBlock struct {
	R			big.Int
	S			big.Int
	PubKey		*ecdsa.PublicKey
	NodeId		uintptr

}

func CreateGenisusBlock(chainId uint, blockNumber unit64, txs string) Block{

	block := Block {chainId, blockNumber, txs, pTxs, pBlockHash}
	return block
}
func CreateBlock(chainId uint, blockNumber unit64, txs string, pTxs string, pBlockHash uint64 ) Block{

	block := Block {chainId, blockNumber, txs, pTxs, pBlockHash}
	return block
}

func (block *Block, prvKey *ecdsa.PrivateKey, nodeId uintptr ) SignBlock() Block{
	kh := NewKeccakState()
	
	data, _ := json.Marshal(block)
	h := HashData(kh , []byte(data))
	
	block.BlockHash := h
	
	r, s, err := Sign(h, prvKey )
	
	pubKey := prvKey.PublicKey()
	signBlock := SignedBlock{r,s,pubKey, nodeId}
	
	block.Signer := signBlock
	
	return block
}
func (block *Block) GetUnsignedBlock()  *Block{

	checkBlock :=Block{}
	checkBlock.ChainId := block.ChainId
	checkBlock.BlockNumber := block.BlockNumber
	checkBlock.Txs := block.Txs
	checkBlock.PTxs := block.PTxs
	checkBlock.PBHash := block.PBHash
	return checkBlock
}

func (block *Block) BlockHash(){
	blockData := block.GetUnsignedBlock()
	jsonData := json.Marshall(blockData)
	kh := common.NewKeccakState()
	block.TxHash = common.HashData(kh, []byte(jsonData))
}

func (block *Block) VerifyBlockHash() bool{
	
	kh := NewKeccakState()
	
	data, _ := json.Marshal(block.GetUnsignedBlock())
	h := HashData(kh , []byte(data))
	return block.BlockHash == h
}

func (block *Block ) VerifyBlockSignature() bool{
	return Verify(block.Hash, block.Signer.R,  block.Signer.S, block.PubKey)

}

func SaveBlockToDisk(block Block){


}

func LoadBlockFromDisk(blockNumber uint64){


}





