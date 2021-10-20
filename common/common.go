package common

import (
	 "bytes"
	 "math/big"
	 "path/filepath"
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
	PB					block.Block						//Last Know Verified Block
	Tx					[]transaction.Transaction		//Last Know Verified Block Transactions
	PBTx				[]transaction.Transaction		//Previous Block Transactions
	BTx					[]transaction.Transaction		//Used to Store Transactions for Pending Block
	PTx					[]crypto.Hash					//Array of Transaction Hashes for Pending Block
	BlockReward			*big.Int						//Amount of FG in Block Reward paid to Block Writers and Leader
	Chain				chain.Chain						//Current Chain
	Chains				chain.Chains					//All Past Year Chains 
	FGValue				float64							//The Value of 1 FG
	Active				[]node.Node						//All Known Active Nodes Next Block
	Nodes				node.Nodes						//All known Nodes
	Writers				[]string						//Array of Current Block Nodes PublicKey as string Based on Block Hash includes Leader wich is the first node listed
	BTxHash				[]crypto.Hash					//Stores processed transaction debit hashes while Block or Leader Node
	PBTxHash			[]crypto.Hash					//Stores previous Block Transactions to account for Transactions sent to Block Leader until block is created & used to validate transactions are in block
	NumTx				int64							//Keeps track of number of Transactions resets at 1,000 Transactions and FGValue is bumped .01
	TTx					[]trasaction.Transaction	    //Used to Transfer Transactions To Nodes One Block at a Time	
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

func CreateBlock(node *node.Node, blockNumber uint64) block.Block{
	NumTxs := uint64(len(PTx))
	var blockTx  []crypto.Hash
	TxFees :=big.NewInt(0)
	if NumTxs > 1000{
		
		for x:=0; x < 1000; x +=1{
			blockTx = append (blockTx, PTx[x])
			percentage := big.NewInt(100)
			for y:=0; y< len(BTx[x].Credit); y +=1{
				txFee := new(big.Int).Div(BTx[x].Credit[y].Amount, percentage)
				TxFees = TxFees.Add(TxFees,txFee)
			}
		}
		go trimPTx()
	}else{
		blockTx = PTx
	}
	
	bn := int64(len(blockTx))
	NumTx += bn
	if bn <1000{
		n:= big.NewInt(0)
		n.SetString("10000000000000000", 10)
		t:= big.NewInt(bn)
		BlockReward :=new(big.Int).Mul(t, n)
	}else{
		BlockReward:= big.NewInt(0)
		BlockReward.SetString("10000000000000000000", 10)
	}
	
	if NumTx > 1000 {
		if FGValue <1000{
			FGValue +=.01 
			
		}
		if (FGValue >=1000) &&(FGValue < 10000){
			FGValue +=.001
			s:=big.NewInt(2)
			BlockReward = new(big.Int).Div(BlockReward, s)

		}
		if FGValue >=10000{
			if FGValue <100000{
				FGValue +=.0001
				b:=big.NewInt(10)
				BlockReward = new(big.Int).Div(BlockReward, b)
			}else{
				b:=big.NewInt(50)
				BlockReward = new(big.Int).Div(BlockReward, b)
			
			}
		}
		
		NumTx = NumTx - 1000
		
	}
	
	NodeTx := PayOutNodes(TxFees, blockNumber)	
	for x:=0; x < len(NodeTx); x +=1{
			blockTx = append (blockTx, NodeTx[x].TxHash)
			BTx	 = append (BTx, NodeTx[x])
			
		}
	WritersTx := PayOutWriters(BlockReward, blockNumber)	
	for x:=0; x < len(WritersTx); x +=1{
			blockTx = append (blockTx, WritersTx[x].TxHash)
			BTx	 = append (BTx, WritersTx[x])
		}
		var block  block.Block
		block.ChainYear = ChainYear
		block.BlockNumber = blockNumber
		block.FGValue = FGValue
		block.Txs = blockTx
		block.NumTxs = NumTxs
		block.Nodes = ActiveNodes
		block.PBHash = PB.BlockHash
		block.NodePayout = NodeTx[0].Debit[0].Amount
		block.WriterPayout = WritersTx[0].Debit[0].Amount
		block.BlockHash = block.HashBlock()
		nodeVals := ElectNodes(block)
	for x:=0; x < len(nodeVals); x +=1{
		block.Writers = append(block.Writers, block.Nodes[nodeVals[x]])
	}

	 return block
}

func ElectNodes(block block.Block) []uint64{
kh :=crypto.NewKeccakState()
	numTx := (PB.NumTxs/500)+1
	var blockNode []uint64
	var numNodes uint64
	numNodes = uint64(len(block.Nodes))
	if numNodes < 7{
		numTx := numNodes
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

func VerifyBlock(block *block.Block){
	numNodes := 0
	for x:=0; x < len(block.Signed); x +=1{
		if block.VerifyWriters(){
				if (len(block.Txs) ==1000) && (bytes.Compare(byte(block.Writers),byte(PB.Writers)) ==0){
					if block.FGValue - PB.FGValue <=.1{
						numNodes +=1
					}
				}else{
					
					if block.Writers ==ElectNodes(block){
						numNodes +=1
					}
				}
			
		}
	}
		//Verify that Block Leader and a Block Node signed the Block
	if numNodes > 1{
		//TODO Store Block to file, Replace Previous Block wtih This Block, 
		block.SaveBlock()
		
	}

}

func PayOutNodes(TxFees *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(int64(len(ActiveNodes)))
payOut := TxFees.Div(TxFees, nodes)

var Txd []transaction.Transaction
	for k,x := range ActiveNodes{
		Tx:=append(Tx, CreatePayoutTransaction(payout, ActiveNodes[x].Id, blockNumber))
	}
	
	return Tx

	
}

func PayOutWriters(blockReward *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(len(Writers))
payOut := blockReward.Div(blockReward, nodes)

var Txd []Transaction
	for x:=0; x<len(Writers); x +=1{
		Tx:=append(Tx, node.CreatePayoutTransaction(payout, Writers[x], blockNumber))
	}
	
	return Tx

	
}


func trimPTx(){
	if len(PTx) > 1000{
		tmpTx := []Transaction
		for x :=1000; x < len(PTx); x+=1{
			tmpTx = append(tmpTx, PTx[x])
		}
		PTx = tmpTx
	}
}


func AllItemsInDir() {
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }

	dir :=filepath.Join(dirname, "fg", "items")
   filepath.Walk(dir, func(path string, info os.FileInfo, e error) {
              if e != nil {
                      fmt.Println(e)
              }

              // check if it is a regular file (not dir)
              if info.Mode().IsRegular() {
                      fmt.Println("file name:", info.Name())
                      fmt.Println("file path:", path)
					  Items[info.Name()] = item.ImportItem(info.Name)
              }
             
      })
	  
}

func genRsa() {
 prvKey := GenerateRSAKey()
 pubKey := prvKey.PublicKey
 secret := RSAEncrypt("Secret", pubKey)
 fmt.Println("Encrypted Message =", secret)
 clearText := RSADecrypt(secret, prvKey)
 fmt.Println("Message =", clearText)
 publicKey := EncodeRSAPubKey(&pubKey)
 fmt.Println("Publick Key =", publicKey)
 pKey := DecodeRSAPubKey(publicKey)
 secret2 := RSAEncrypt("Secret2", pKey)
 fmt.Println("Encrypted Message =", secret2)
 clearText2 := RSADecrypt(secret2, prvKey)
 fmt.Println("Message =", clearText2)
 pk := prvKey
 err:= StoreRSAKey( pk ,"Pass", "Key1")
 fmt.Println(err)
 pvKey, pbKey, err :=GetRSAKey("Key1", "Pass") //rsa.PrivateKey,rsa.PublicKey, error
 secret3 := RSAEncrypt("Secret3", pbKey)
 fmt.Println("Encrypted Message =", secret3)
 clearText3 := RSADecrypt(secret3, pvKey)
 fmt.Println("Message =", clearText3)


}