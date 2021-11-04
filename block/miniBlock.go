package block

import(
		"encoding/json"
		"fmt"
		"github.com/fgeth/fg/crypto"
)


type MinBlock struct{
	ChainYear			uint64						//Year this block was created
	BlockNumber			uint64						//Block Number
}

func (block MinBlock) BlockHash() string{
		
	json , err:= json.Marshal(block)
	if err !=nil{
		fmt.Println("Error reading MiniBLock", err)
		
	}
	
	blockHash :=crypto.HashTx([]byte(json))
	fmt.Println("The Hash :", blockHash)
	return blockHash
}