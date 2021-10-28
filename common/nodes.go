package common

import(
	"bytes"
	 "fmt" 
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/transaction"
)


func CompareWriters(writers []string, theWriters []string) bool{
	NotEqual := false
	for x :=0; x < len(theWriters); x +=1{
		if writers[x] != theWriters[x]{
			NotEqual = true
		}
			
	}
	return NotEqual
}

func ElectNodes(block block.Block) []uint64{

	numTx := (PB.NumTxs/500)+1
	var blockNode []uint64
	var numNodes uint64
	numNodes = uint64(len(block.Nodes))
	if numNodes < 7{
		numTx = numNodes
	}else{
		if numTx < 7 {
			numTx =7
		}else{
			if numTx > numNodes{
				numTx = numNodes
			}
		}
	}
	if bytes.Compare(block.BlockHash, block.HashBlock()) ==0 {
	    uintA, uintB, uintC, uintD :=crypto.HashToUint64(block.BlockHash)
		
		tmp := uint64(0)
		bn := uint64(0)
		for x :=uint64(0); x < numTx; x +=1{
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
			
			
			
	}
	return blockNode
}







func SelectNode(tx transaction.Transaction) []uint64{

	kh :=crypto.NewKeccakState()
	txData := tx.TxData()
	txhash := crypto.HashData(kh, []byte(fmt.Sprintf("%v", txData)))
	uintA, uintB, uintC, uintD :=crypto.HashToUint64(txhash)
		var blockNode []uint64
		tmp := uint64(0)
		bn := uint64(0)
		numNodes := uint64(len(Writers))
		for x :=uint64(0); x < numNodes; x +=1{
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