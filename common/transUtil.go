package common

import(
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt" 
	"io/ioutil"
	"math/big"
	"net/http"
	"path/filepath"
	//"strconv"
	"time"
	"github.com/fgeth/fg/bank"
	"github.com/fgeth/fg/block"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/note"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/transaction"
	"github.com/fgeth/fg/wallet"
	
	
)


func PayOutNodes(TxFees *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(int64(len(ActiveNodes)))
payOut := TxFees.Div(TxFees, nodes)

var Txd []transaction.Transaction
	for k,x := range ActiveNodes{
		Txd =append(Tx, CreatePayoutTransaction(payOut, x, blockNumber))
		fmt.Println("Paying node number:", k)
	}
	
	return Txd

	
}

func PayOutWriters(blockReward *big.Int, blockNumber uint64)[]transaction.Transaction{

nodes := big.NewInt(int64(len(Writers)))
payOut := blockReward.Div(blockReward, nodes)

var Txd []transaction.Transaction
	for x:=0; x<len(Writers); x +=1{
		Txd=append(Tx, CreatePayoutTransaction(payOut, Writers[x], blockNumber))
	}
	
	return Txd

	
}

func CreatePayoutTransaction(amt *big.Int, pubKey string, blockNumber uint64) transaction.Transaction{


	var Debit transaction.BaseTransaction
	var Credit transaction.BaseTransaction
	var Tx transaction.Transaction
	Debit.ChainYear = ChainYear
	Debit.BlockNumber = blockNumber
	Debit.Time = time.Now()
	Debit.Amount = amt
	Debit.TxHash = Debit.HashBaseTx(pubKey)
	Credit.ChainYear = ChainYear
	Credit.BlockNumber = blockNumber
	Credit.Time = time.Now()
	Credit.TxHash = Credit.HashBaseTx(pubKey)
	Credit.Amount = amt
	Tx.Change.TxHash = Tx.Change.HashBaseTx(pubKey)
	Tx.Debit =Debit
	Tx.Credit = append(Tx.Credit, Credit)
	Tx.TxHash = Tx.HashTx()
	Tx.Credit[0].R, Tx.Credit[0].S = crypto.TxSign([]byte(Tx.TxHash), MyNode.PrvKey)
	Tx.Payout = true
	return Tx
	
}

func CreateTransaction(amt *big.Int, Credit []transaction.BaseTransaction, pubKey string, pbKey string, blockNumber uint64, PrvKeys []*ecdsa.PrivateKey) transaction.Transaction{


	
	var Tx transaction.Transaction
	
	Tx.Change.TxHash = Tx.Change.HashBaseTx(pbKey)
	Tx.Debit = CreateDebitTxs(amt, pubKey, blockNumber)
	for x:=0; x < len(Credit); x+=1{
		Tx.Credit = append(Tx.Credit, Credit[x])
	}
	fmt.Println("TxHash :", Tx.HashTx())
	Tx.TxHash = Tx.HashTx()
	Tx.Debit.TxId = Tx.TxHash
	for x:=0; x < len(Tx.Credit); x +=1{
		pubKey = crypto.EncodePubKey(&PrvKeys[x].PublicKey)
		Tx.Credit[x].R, Tx.Credit[x].S = crypto.TxSign([]byte(Tx.TxHash), PrvKeys[x])
		Tx.Credit[x].OTP = pubKey
		Tx.Credit[x].Spent = Tx.TxHash
	}
	Tx.Payout = true
	return Tx
	
}

func CreateDebitTxs(amt *big.Int, pubKey string, blockNumber uint64) transaction.BaseTransaction{
	var Debit transaction.BaseTransaction
	Debit.Amount = amt
	Debit.ChainYear = ChainYear
	Debit.BlockNumber = blockNumber
	Debit.Time = time.Now()
	Debit.TxHash = Debit.HashBaseTx(pubKey)
	fmt.Println("Debit.HashBaseTx :", Debit.HashBaseTx(pubKey))
	return Debit
}


//TODO Get Txs
func GetTxs(){

}

func VaildTransaction(Tx transaction.Transaction) bool{
txFee := Tx.CalcFee()
 txInterest := big.NewInt(0)
if CoinValue >=100{
	txInterest = Tx.CalcInterest()
}
if Tx.Payout == true{
			if Tx.Credit[0].VerifySig(){
				block := block.ImportBlock(Tx.Credit[0].ChainYear , Tx.Credit[0].BlockNumber, MyNode.Path)
				for x:=0; x<len(block.Writers); x +=1{
					if block.Writers[x] == Tx.Credit[0].OTP{
						if Tx.Debit.Amount == block.NodePayout{
								return true
							}
						if Tx.Debit.Amount == block.WriterPayout{
							return true
						}
					}
				}
			}
		
		}

				TC := Tx.Credits()
				TD := Tx.Debit.Amount
				ValidCredit :=0
				if TC.Add(TC,txInterest) == txFee.Add(txFee, TD.Add(TD, Tx.Change.Amount)){
				for x:=0; x< len(Tx.Credit); x+=1{
					cTx := transaction.ImportBaseTx([]byte(Tx.Credit[x].TxHash), MyNode.Path)
					if cTx.OTP ==""{
						cTx.OTP=Tx.Credit[x].OTP
						if cTx.VerifySig(){
							cTx.SaveTx(MyNode.Path)
							ValidCredit +=1
						}else{
							return false  //not signed
						}
						
					}else{
						return false		//Double Spend
					}
				}
				if ValidCredit >= len(Tx.Credit){
					Tx.Debit.SaveTx(MyNode.Path)
					Tx.Change.SaveTx(MyNode.Path)
					return true
				}
				
					
			}
return false
}



//TODO fix where this imports the transactions
func ImportTx(txHash string) {
	//dirname, err := os.UserHomeDir()
    //if err != nil {
    //   fmt.Println( err )
    //}
	
	path :=filepath.Join(MyNode.Path, "tx", txHash )
	file, _ := ioutil.ReadFile(path)
	var tx transaction.Transaction
	_ = json.Unmarshal([]byte(file), &tx)
}



func SendTransaction(tx transaction.Transaction) bool{
	blockNodes := SelectNode(tx)
	for x :=0; x<len(blockNodes); x+=1{
		if SubmitTransaction(tx, Writers[blockNodes[x]]){
			return true
		}
	}
	return false
}

func SubmitTransaction(tx transaction.Transaction, writer string) bool{
	json, _:= json.Marshal(tx)
	
	 if node, ok :=TheNodes.Node[writer]; ok{
		call := "sendTx"
		//call = block, node, tx, or account
		url := fmt.Sprintf("http://%i:%p/%t", node.Ip, node.Port, call)

			//req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
			if err != nil {
			  fmt.Println("Error reading response. ", err)
			
			  return false
			}
			
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body. ", err)
			}

			fmt.Printf("%s\n", body)
			
			return true
		
	}
	return false

}

func SubmitPayment(payment note.Stack, theItem item.Buy) bool{
	var response bool
	response = false
	numB := len(Ring.Banks)
	if numB>0{
		b :=findBank(uint64(numB),payment.Id )
		
		p, err:= json.Marshal(payment)
		if err ==nil{
			i, err:= json.Marshal(theItem)
			if err ==nil{
				for x :=0; x < len(b); x +=1{
					securePayment, errP :=crypto.RSAEncryptByte(p, Ring.Banks[b[x]].PubKey) 
				
				
					secureItem, errI :=crypto.RSAEncryptByte(i, Ring.Banks[b[x]].PubKey) 
					if errP ==nil && errI == nil{
						BM := bank.BankMessage {Ring.Banks[b[x]].Id, securePayment, secureItem}
						json, _ := json.Marshal(BM)
						url :="http://"+Ring.Table[0].Node.Ip+":"+Ring.Table[0].Node.Port+"/bankMessage"
						fmt.Println("Connecting to Ring at :", url)
						_, err := http.Post(url, "application/json", bytes.NewBuffer(json))
						if err ==nil{
							response = true
							break
							
						}
					}
				}
			}
		}
	}
	return response

}

func SendFunds(wallet wallet.Wallet, theItem item.Buy) string{
  if wallet.Dollars < theItem.Amount{
	fmt.Println("Wallet has : ", wallet.Dollars)
	fmt.Println("Not enough Funds")
	fmt.Println("Item is :", theItem.Amount)
	return "Not enough funds"
  }
  fmt.Println("Item is :", theItem.Amount)
  fmt.Println("Wallet has : ", wallet.Dollars)
  
  var aStack note.Stack
  var amount float64

  amount =float64(0)
	for k,v := range wallet.Notes{

			if amount < theItem.Amount{
				aStack.Notes = map[string]note.Note{v.Id:v}
				amount += Coins2VDollars(v.Coins)
				
				 delete(wallet.Notes, k)

			}
			if amount >= theItem.Amount{
				break
			}
		

	}
	aStack.Amount = amount

	response:= SubmitPayment(aStack, theItem)
	if !response{
		return "Error Sending Payment"
		 for k,n :=range aStack.Notes {
			wallet.Notes[k] = n
		 }
	}
	
	return "Sent Payment"
	

}

func findBank(banks uint64, id string )[]int{
var tmp, bn uint64
var theBanks []int
uintA, uintB, uintC, uintD :=crypto.B32HashToUint64([]byte(id))
  				for x :=uint64(0); x < banks; x +=uint64(1){
						switch x%4{
							case  1:
								tmp = (uintA % banks)+x
								
									
							case 2:
								tmp = (uintB % banks)+x
							
							case 3:
								tmp = (uintC % banks)+x
							
							case 4:
								tmp = (uintD % banks)+x
				
						}	
							if tmp > banks{
								bn = tmp - banks
							}else{
								bn = tmp
							}
							theBanks  = append(theBanks, int(bn))
						
					}
	
		return theBanks
}