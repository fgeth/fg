
package chain

import(

)


type Chain struct {
	ChainYear		uint32
	BlockNumber		uint64				//Highest Block For This Chain

}

type Chains struct {
	Chain 		map[uint32]Chain

}