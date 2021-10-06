package block

import (
	"crypto/ecdsa"
    "crypto/sha1"
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"bufio"
	"hash"
	"math/big"
	"sync"
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	 "os"
	 "github.com/fgeth/fg/crypto"
)

type Block struct {
	BlockNumber			uint64
	Txs					[]transaction	   		//Array with completed transactions
	PTxs				[]transaction			//Array of pending transactions
	Nodes				string					//Comma seperated list of all New Nodes IPS on network in Node order
	RmNodes				string					//Comma seperated list of all nodes that did not respond during this Block Nodes start at its Node Array index and goes through array
	PBHash				uint64					//Hash of previous Block
	NBLNode				string					//IP of the Next Block Node Leader based on Block Hash
	NBNodes				string					//Comma Seperated list of the remaining Block Nodes IPS  Based on Block Hash
	BlockHash			uint64					//Hash of this Block including previous Blocks Hash & list of new Nodes and Nodes to remove
	ChainHash			uint64					//Hash of all the Block hashes up to this point for this chain includes this BlockHash
	Signers				[]SignedTx				//Signature of Block Nodes
}




func CreateBlock(chain Chain) Block{
	hash := getHash(chain.GetNextBlock(), GetTxs(), GetNodes, chain.GetLastBlock)
	block := Block {chain.GetNextBlock(), GetTxs(), chain.GetBlockHash1(), chain.GetBlockHash2(), chain.GetBlockHash3(), hash}
	return block
}

func SaveBlockToDisk(block Block){


}

func LoadBlockFromDisk(blockNumber uint64){


}





