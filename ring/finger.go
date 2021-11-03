package ring

import (

	"github.com/fgeth/fg/node"
	
)

type FingerTable struct {
	Id			uint64			//Location on the ring
	Node		node.SNode		//The Node
	
}


