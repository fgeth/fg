package writer

import(
	
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/transaction"
)

type Writer struct{
	Block 				block.Block						//New Block
	Txs					[]transaction.Transaction		//Array of new Block Transactions
	ActiveNodes			[]string						//Array of new Active Nodes

}