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
	 "github.com/fgeth/fge/crypto"
)

type Block struct {
	BlockNumber			uint64
	Txs					[]uint64       			//Array with transaction hashes
	Nodes				string					//Comma seperated list of all New Nodes IPS on network in Node order
	PBHash				uint64					//Hash of previous Block
	NBLNode				string					//IP of the Next Block Node Leader
	NBNodes				string					//Comma Seperated list of the remaining Block Nodes IPS
	BlockHash			uint64					//Hash of this Block including previous Blocks Hash, list of Nodes, Next Block Nodes and Block Leader	
	ChainHash			uint64					//Hash of all the Block hashes up to this point for this chain includes this BlockHash
	Signer				SignedTx				//Signature of Block Leader
}




func CreateBlock(chain Chain) Block{
	hash := getHash(chain.GetNextBlock(), GetTxs(), GetNodes, chain.GetLastBlock)
	block := Block {chain.GetNextBlock(), GetTxs(), chain.GetBlockHash1(), chain.GetBlockHash2(), chain.GetBlockHash3(), hash}
	return block
}




