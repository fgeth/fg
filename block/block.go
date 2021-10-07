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
	//Nodes				string						//Comma seperated list of all New Nodes IPS on network in Node order
	//RmNodes			string						//Comma seperated list of all nodes that did not respond during this Block Nodes start at its Node Array index and goes through array
	PBHash				common.Hash					//Hash of previous Block
	NBLNode				untptr						//ID of the Next Block Node Leader based on Block Hash
	Writers				[]untptr					//array of the Block Nodes untptr  Based on Block Hash includes Leader
	BlockHash			common.Hash					//Hash of this Block including previous Blocks Hash & list of new Nodes and Nodes to remove
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
	checkBlock.NBLNode := block.NBLNode
	checkBlock.NBNodes := block.NBNodes
	checkBlock.NBLNode := block.NBLNode
	return checkBlock
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





