package common

import (
	 "github.com/fgeth/fg/block"
	 "github.com/fgeth/fg/chain"
	 "github.com/fgeth/fg/crypto"
	 "github.com/fgeth/fg/item"
	 "github.com/fgeth/fg/node"
	 "github.com/fgeth/fg/transaction"
)

var (
	ChainYear			uint64							//Current Year
	BlockNumber			uint64							//Current Block
	ActiveNodes			[]string						//Array of known active Nodes Public Key as string
	PB					Block							//Last Know Verified Block
	Tx					[]Transaction					//Last Know Verified Block Transactions
	PBTx				[]Transaction					//Previous Block Transactions
	BTx					[]Transaction					//Used to Store Transactions for Pending Block
	PTx					[]Crypto.Hash					//Array of Transaction Hashes for Pending Block
	BlockReward			big.Int							//Amount of FG in Block Reward paid to Block Writers and Leader
	Chain				Chain							//Current Chain
	Chains				Chains							//All Past Year Chains 
	FGValue				float64							//The Value of 1 FG
	Active				[]Node							//All Known Active Nodes Next Block
	Nodes				Nodes							//All known Nodes
	Writers				[]string						//Array of Current Block Nodes PublicKey as string Based on Block Hash includes Leader wich is the first node listed
	BTx					[]Crypto.Hash					//Stores processed transaction debit hashes while Block or Leader Node
	PBTx				[]crypto.Hash					//Stores previous Block Transactions to account for Transactions sent to Block Leader until block is created & used to validate transactions are in block
	NumTx				uint32							//Keeps track of number of Transactions resets at 1,000 Transactions and FGValue is bumped .01
	TTx					[]Transaction					//Used to Transfer Transactions To Nodes One Block at a Time	
	Items				Items							//Index is Item Id
	
)


func FG2USD(amount *big.Int) float64{
    fg := new(big.Int)
	fg.SetString("1000000000000000000", 10)

	f := new(big.Float).SetInt(amount)
	t := new(big.Float).SetInt(fg)
	f = f.Quo(f, t)

	fv, _:= f.Float64()
	
	usd :=   FGValue * fv
	return usd
	
	
}

func USD2FG(amount float64) *big.Int{
	
	bigval := new(big.Float)

	fgs := amount / FGValue

	bigval.SetFloat64(fgs)

	fv := new(big.Float)
	fv.SetString("1000000000000000000")

	fg :=new(big.Int)

	bigval.Mul(bigval, fv)
	bigval.Int(fg)

	return fg
	
	
}

func (node *Node) CreateBlock(blockNumber uint64) Block{
	NumTxs := len(PTx)
	blockTx := []Transaction{}
	
	if NumTxs > 1000{
		TxFees :=big.NewInt(0)
		for x:=0; x < 1000; x +=1{
			blockTx = append (blockTx, PTx[x])
			percentage := big.NewInt(100)
			
			txFee := new(big.Int).Div(PTx[x].Credit.Amount, percentage)
			TxFees = TxFees.Add(TxFees,txFee)
		}
		go trimPTx()
	}else{
		blockTx = PTx
	}
	
	bn : = len(blockTx)
	NumTx += bn
	if bn <1000{
		n:= big.SetString("10000000000000000", 10)
		t:= big.NewInt(bn)
		BlockReward :=new(big.Int).Mul(bn, n)
	}else{
		BlockReward := big.SetString("10000000000000000000", 10)
	}
	
	if NumTx > 1000 {
		if FGValue <1000{
			FGValue +=.01 
			
		}
		if (FGValue =>1000) &(FGValue < 10000){
			FGValue +=.001
			n=big.NewInt(2)
			BlockReward = new(big.Int).div(BlockReward, n)

		}
		if FGValue =>10000{
			if FGValue <100000{
				FGValue +=.0001
				n=big.NewInt(10)
				BlockReward = new(big.Int).div(BlockReward, n)
			}else{
				n=big.NewInt(50)
				BlockReward = new(big.Int).div(BlockReward, n)
			
			}
		}
		
		NumTx = NumTx - 1000
		
	}
	
	NodeTx := PayOutNodes(TxFees, blockNumber)	
	for x:=0; x < len(NodeTx); x +=1{
			blockTx = append (blockTx, NodeTx[x])
		}
	WritersTx := PayOutWriters(BlockReward, blockNumber)	
	for x:=0; x < len(WritersTx); x +=1{
			blockTx = append (blockTx, WritersTx[x])
		}
	block = Block{ChainYear, blockNumber, FGValue, blockTx, ActiveNodes, NumTxs, PB.BlockHash}
	block.NodePayOut = NodeTx[0].Debit.Amount
	block.WriterPayOut = WritersTx[0].Debit.Amount
	block.BlockHash = block.HashBlock()
	nodeVals := block.ElectNodes()
	for x:=0; x < len(nodeVals); x +=1{
		block.Writers = append(block.Writers, block.Nodes[nodeVals[x]])
	}

	 return block
}


func PayOutNodes(TxFees big.Int, blockNumber uint64)[]Transaction{

nodes := big.NewInt(len(Nodes))
payOut := TxFees.Div(TxFees, nodes)

var Txd []Transaction
	for k,x := range Nodes{
		Tx:=append(Tx, node.CreatePayoutTransaction(payout, Nodes[x].Id, blockNumber))
	}
	
	return Tx

	
}

func PayOutWriters(blockReward big.Int, blockNumber uint64)[]Transaction{

nodes := big.NewInt(len(Writers))
payOut := blockReward.Div(blockReward, nodes)

var Txd []Transaction
	for x:=0; x<len(Writers); x +=1{
		Tx:=append(Tx, node.CreatePayoutTransaction(payout, Writers[x], blockNumber))
	}
	
	return Tx

	
}