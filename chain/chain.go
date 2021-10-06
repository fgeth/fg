package chain

import (
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
	 "github.com/fgeth/fge/block"
	 
)

type Chain struct {
	ChainId			uint					//Current Year -- New Chain is started each year
	BlockNumber		uint64					//Current Years' Current Block Number
	Blocks			[]Block					//Array of the curret years' Blocks
	Writers			[]unitptr				//Current list of Nodes that have write access to the chain list of the current BlockNodes and Pervious BlockNodes Writer[0] points to BlockLeaderNode
	LYLBlockNum		common.Hash				//Last Years Last Block Number
	LYLBlockHash	common.Hash				//Last Years Last Block Hash
	Hash			[]common.Hash	        //Array of Hashes when a new block is added to the chain all block hashes are put together and hashed and the result is 
	PYH				[]common.Hash			//appended to this array so that the last hash is the hash for the current height of the BlockChain
											//PYH is an array of the all previous years final hashes
}


func NewChain(year uint, LYBN uint64, LYLBH unit64, writers []untptr, hash []uint64, PYH [] uint64) Chain{
	blockNumber := uint64(0)
	blocks := make([]Block, 2,420,069)   // 2,102,400 potential blocks at one every 15 seconds
	chain := new Chain{year, blockNumber, blocks, writers, LYBN, LYBH, hash, PYH}

	return chain
}
func (chain *Chain)GetNextBlock(){
	return BlockNumber +1
}

func FirstChain(year uint, writers unitptr) Chain{

	blocks := make([]Block, 2,420,069)   // 2,102,400 potential blocks at one every 15 seconds
	blockNumber := uint64(0)
	var LYBN := common.Hash{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	var LYBH :=common.Hash{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	fHash : []common.Hash{LYBN, LYBH}
	pYH : []common.Hash{LYBN, LYBH}
	chain := new Chain{year, blockNumber, blocks, writers, LYBN, LYBH, fHash, pYH}

	return chain
}
func (chain *Chain)GetNextBlock(){
	return BlockNumber +1
}



func (chain *Chain) GetBlockHash1(){
	return chain.hash[chain.BlockNumber]
}

func(chain *Chain) GetBlockHash2(){
    return chain.hash[chain.BlockNumber-1]
}
func (chain *Chain)GetBlockHash3(){
	eturn chain.hash[chain.BlockNumber-2]
}

func SaveChainToDisk(chain Chain){


}

func LoadChainFromDisk(chainId uint){


}


