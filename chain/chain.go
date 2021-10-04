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
	LYLBlockNum		unit64					//Last Years Last Block Number
	LYLBlockHash	unit64					//Last Years Last Block Hash
	Writer			[]unitptr				//Current list of Nodes that have write access to the chain list of the current BlockNodes and Pervious BlockNodes Writer[0] points to BlockLeaderNode
	Hash			[]unit64		        //Array of Hashes when a new block is added to the chain all block hashes are put together and hashed and the result is 
	PYH				[]unit64				//appended to this array so that the last hash is the hash for the current height of the BlockChain
											//PYH is an array of the all previous years final hashes
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