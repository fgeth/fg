
package chain

import(
 	"github.com/fgeth/fg/block"
)


type Chain struct {
	ChainYear		uint32
	BlockNumber		uint64				//Highest Block For This Chain
	Blocks			[]block.Block

}

type Chains struct {
	Chain 		map[uint32]Chain		//Index is chain year

}