package common

import(
     "bytes"
	 "fmt" 
	"encoding/json"
	"io/ioutil"
	"math/big"
	//"net/http"
	//"net/url"
	"os"
	"path/filepath"
	 "strconv"
	 "time"
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fasthttp"
	"github.com/fgeth/fasthttp/fasthttpproxy"
)


//TODO Fix To Where this imports the blocks
func ImportBlocks(blockNumber uint64) {
	//dirname, err := os.UserHomeDir()
    //if err != nil {
     //   fmt.Println( err )
    //}
   // fmt.Println( dirname )
	path :=filepath.Join(MyNode.Path, "chain", strconv.FormatUint(ChainYear, 10))
	for x:=uint64(0); x < blockNumber; x +=uint64(1){
		fileName := filepath.Join(path, strconv.FormatUint(x, 10))
		var block block.Block
		file, err := ioutil.ReadFile(fileName)
		if err == nil{
			
			_ = json.Unmarshal([]byte(file), &block)
		}else{
			block = GetBlock(x, MyNode.OA)
		}
		Chain.Blocks = append(Chain.Blocks, block)
	}

}
//TODO Sign Genesis block
func SignGenesisBlocks(){

}

func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}
//Accepts blockNumber and Onion Address to get block
func GetBlock(x uint64, OA string) block.Block{
		//var block1 block.MinBlock
		var block2 block.Block
		//block1.BlockNumber = x
		//block1.ChainYear = ChainYear
		//data, err:= json.Marshal(block1)
		var dst []byte
		BN :=fasthttp.ArgsKV{[]byte("BlockNumber"), i64tob(x), false}
		CY :=fasthttp.ArgsKV{[]byte("ChainYear"), i64tob(ChainYear), false}
		var postArgs fasthttp.Args
		postArgs.Args = append(postArgs.Args, BN)
		postArgs.Args = append(postArgs.Args, CY)
		
		
		//type Args struct {
		//noCopy noCopy //nolint:unused,structcheck
		//	args []argsKV
		//	buf  []byte
		//}

		//type argsKV struct {
		//key     []byte
		//value   []byte
		//noValue bool
		//}
		//if err !=nil{
		//	fmt.Println("Error Marshalling Block")
		//}
		//fmt.Println("Data :", data)
		call := "getBlock"
		//call = block, node, tx, or account
		url1 := OA +"/"+call
		fmt.Println("url:", url1)
		err := TorDialer(url1)
		if err ==nil{
				//bytes.NewBuffer
			//p := "http://127.0.0.1:"+MyNode.Tor	
			//proxy, _ := url.Parse(p) 
			//http.DefaultTransport = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}	
			c := &fasthttp.Client{
				Dial: fasthttpproxy.FasthttpSocksDialer("socks5://localhost:9050"),
			}
			//statusCode int, body []byte, err error
			status, resp, err := c.Post( dst, url1, &postArgs)

			if err != nil {
			  // Error reading Block data
			  fmt.Println("Error reading Block", err)
			}else{
			
			//defer resp.Body.Close()

			fmt.Println("response Status:", status)
			//fmt.Println("response Headers:", resp.Header)

			//body, err := ioutil.ReadAll(resp)
			//if err != nil {
			//	fmt.Println("Error reading body. ", err)
			//}

			//fmt.Printf("%s\n", body)
			 json.Unmarshal(resp, &block2)
			}
		}

	return block2
}
//TODO GetBlocks
func GetBlocks(){

}

func CreateBlock( ) block.Block{
var block  block.Block
if MyNode.Leader{
	blockNumber:= BlockNumber + uint64(1)
	NumTxs := uint64(len(PTx))
	var blockTx  []string
	TxFees :=big.NewInt(0)
	MultiBlock :=false
	if NumTxs > 1000{
		MultiBlock = true
		for x:=0; x < 1000; x +=1{
			blockTx = append (blockTx, PTx[x])
			percentage := big.NewInt(500)
			for y:=0; y< len(BTx[x].Credit); y +=1{
				txFee := new(big.Int).Div(BTx[x].Credit[y].Amount, percentage)
				TxFees = TxFees.Add(TxFees,txFee)
			}
		}
		
	}else{
		blockTx = PTx
	}
	BlockReward :=big.NewInt(0)
	bn := int64(len(blockTx))
	NumTx += bn
	if bn <1000{
		n:= big.NewInt(0)
		n.SetString("10000000000000000", 10)
		t:= big.NewInt(bn)
		BlockReward =new(big.Int).Mul(t, n)
	}else{
		BlockReward= big.NewInt(0)
		BlockReward.SetString("10000000000000000000", 10)
	}
		fmt.Println(BlockReward)
	
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
			AddBTX(NodeTx[x])
			
			
		}
	WritersTx := PayOutWriters(BlockReward, blockNumber)	
	for x:=0; x < len(WritersTx); x +=1{
			blockTx = append (blockTx, WritersTx[x].TxHash)
			AddBTX(WritersTx[x])
		}
		
		t := time.Now()
		if uint64(t.Year())> ChainYear{
			IncChainYear()
			block.ChainYear = ChainYear
		}
		block.ChainYear = ChainYear
		block.BlockNumber = blockNumber
		block.FGValue = FGValue
		block.Txs = blockTx
		block.NumTxs = NumTxs
		block.Nodes = ActiveNodes
		block.PBHash = PB.BlockHash
		block.NodePayout = NodeTx[0].Debit.Amount
		block.WriterPayout = WritersTx[0].Debit.Amount
		block.BlockHash = block.HashBlock()
		if MultiBlock{
			block.Writers = PB.Writers
			go trimPTx()
		}else{
			nodeVals := ElectNodes(block)
			for x:=0; x < len(nodeVals); x +=1{
				block.Writers = append(block.Writers, block.Nodes[nodeVals[x]])
			}
		}
			
		SwapBlocks(&block)
		block.SaveBlock(MyNode.Path)
	 
	}
	return block
}


//TODO Recreate Block Failed
func BlockFailed(blockNumber uint64){

}


func VerifyBlock(block *block.Block) bool{
	numNodes := 0
	for x:=0; x < len(block.Signed); x +=1{
		if block.VerifyBlock(PB){
				if (len(block.Txs) ==1000) && (CompareWriters(block.Writers,PB.Writers)){
					if block.FGValue - PB.FGValue <=.01{
						numNodes +=1
					}
				}else{
					
					if CompareWriters(block.Writers, block.GetWriters(ElectNodes(*block))){
						numNodes +=1
					}
				}
			
		}
	}
		//Verify that Block Leader and a Block Node signed the Block
	if numNodes > 1{
		//TODO Store Block to file, Replace Previous Block wtih This Block, 
		if bytes.Compare(block.BlockHash, block.HashBlock()) ==0{	
			
			if bytes.Compare(block.PBHash, PB.BlockHash) ==0{
				if CheckBlock(block){
					CB := ImportBlock(block)		// current saved block need to add the other block signatures
					for x:=0; x < len(CB.Signed); x +=1{
						if CB.VerifyBlock(PB){
							block.Signed = append(block.Signed, CB.Signed[x])
						}
					}
				}
				block.SaveBlock(MyNode.Path)
				if BlockNumber < block.BlockNumber{
					SwapBlocks(block)
				}
				return true
			}
		}
		
	}
	return false

}

func CheckBlock(block *block.Block) bool{
	dirname, err := os.UserHomeDir()

    if err != nil {
        fmt.Println( err )
    }
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "fg", "chain", strconv.FormatUint(block.ChainYear, 10))

	fileName := filepath.Join(path, strconv.FormatUint(block.BlockNumber, 10))
	myfile, e := os.Stat(fileName)
	if e != nil{
	  fmt.Println( e )
	  return false
	}
	fmt.Println( myfile )
	return true

}

func ImportBlock(CB *block.Block) block.Block{
	return block.ImportBlock(CB.ChainYear, CB.BlockNumber, MyNode.Path)

}