package wallet


import(
	"crypto/ecdsa"
	"math/big"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/fgeth/fg/bank"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/note"
	"github.com/fgeth/fg/transaction"

)

type Wallet struct {
	Id			string							//Address
	FGs			float64
	Wei			*big.Int
	Dollars		float64
	FGValue		float64
	Items		Selling							//Items that are for sell
	Buy			Buying							//Items we have bought
	Notes		[]note.Note						//Index is 100, 20, 10, 5, 1, or 0	
	Stack		[]note.Stack
	Debits		Debits
	Auth		string
}


type Selling struct {
	Item			map[string]item.Item								 //Index of Item Id and the Item
	Keys			map[string][]*ecdsa.PrivateKey						//Index is Item Id and array of private keys for the transaction
	
}

type Buying struct {
	Item			map[string]item.Item								 //Index is Item Id and the Item
	Tx				map[string][]transaction.Transaction				//Index is Item Id	and the array of transactions that go with that Item
	Keys			map[string][]*ecdsa.PrivateKey						//Index is Item Id and array of private keys for the transaction

}

type Debits struct {
	Debit	map[string]transaction.BaseTransaction								//Index of debit Tranasaction Hash with transaction
	
}




func (wallet Wallet) SaveWallet(dirname string){

	path :=filepath.Join(dirname, "wallet")
	 fmt.Println("Path ", path)
	_, err := os.Stat(dirname)
    if err !=nil {
		fmt.Println("Wallet Root directory does not Exist", err)
		err := os.Mkdir(dirname, 0755)
		if err !=nil{
			fmt.Println("Failed to make Wallet root directory", err)
			dirname, _ := os.UserHomeDir()
			
			path =filepath.Join(dirname, "wallet")
			_, err = os.Stat(path)
			if err !=nil{
				err = os.Mkdir(dirname, 0755)
				if err !=nil{
					fmt.Println("failed to make root directory", err)
				}
				err = os.Mkdir(path, 0755)
				if err !=nil{
					fmt.Println("failed to make wallet directory", err)
				}
			}
		}else{
			err = os.Mkdir(filepath.Join(dirname, "wallet"), 0755)
			if err !=nil{
				fmt.Println("failed to make wallet directory", err)
			}
		}

    }else{
		_, err := os.Stat(path)
			if err !=nil {
				err = os.Mkdir(filepath.Join(dirname, "wallet"), 0755)
				if err !=nil{
					fmt.Println("failed to make wallet directory", err)
				}
			}
	}
  
	fileName := filepath.Join(path, wallet.Id)

	file, err := json.MarshalIndent(wallet, "", " ")
	if err !=nil{
		fmt.Println("Error Marshalling Wallet :", err)
	}
	fmt.Println("The wallet Marshalled: ", file)
		file, _= crypto.Encrypt([]byte(wallet.Auth), file)
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}


func ImportWallet(dirname, walletId, auth string) (Wallet, error){
	
	
	var wallet Wallet
    var errWal Wallet

	path :=filepath.Join(dirname, "wallet")

	fileName := filepath.Join(path, walletId)
	fmt.Println("File Name : ", fileName )
	_, e := os.Stat(fileName)
	if e != nil{
		dirname, _ := os.UserHomeDir()
		path :=filepath.Join(dirname, "wallet")
		fileName := filepath.Join(path, walletId)
		//fmt.Println("File Name : ", fileName )
		_, e1 := os.Stat(fileName)
		
		if e1 != nil{
			return errWal, e1
			
		}else{
			file, err:= ioutil.ReadFile(fileName)
			if err == nil{
				//fmt.Println("Unmarshalling File : ", fileName )
				theFile, _:= crypto.Decrypt([]byte(auth), file)
				err :=json.Unmarshal([]byte(theFile), &wallet)
				
				if err != nil {
					fmt.Println("couldn't unmarshal parameters", err)
					return errWal, err

				}
			}else{
				return errWal, err
			}
		}
		
		//fmt.Println( e )
	}else{
		file, err:= ioutil.ReadFile(fileName)
		if err ==nil{
			theFile,_ := crypto.Decrypt([]byte(auth), file)
			err := json.Unmarshal([]byte(theFile), &wallet)
			
			fmt.Println("Unmarshalling File : ", fileName )
			if err !=nil{
				fmt.Println("Error Unmarshalling File : ", err )
			}
		}
	if err != nil {
        fmt.Println("couldn't unmarshal parameters", err)
			return errWal, err
    }
	}


	return wallet, nil
	
}


func(wallet Wallet) SendFunds(theItem item.Buy) string{
  if wallet.Dollars < theItem.Amount{
	return "Not enough Funds"
  }
  var aStack note.Stack
  var amount float64
  var tmpNotes []note.Note
  amount =float64(0)
	for x:= len(wallet.Notes)-1; x >0; x-=1{
		if amount < theItem.Amount{
			aStack.Notes = append(aStack.Notes, wallet.Notes[x])
			amount += wallet.Notes[x].Amount
			aStack.Amount = amount
			tmpNotes = removeNote(wallet.Notes)

		}
		if amount >= theItem.Amount{
			break
		}

	}

	change, err:= bank.SubmitPayment(aStack, theItem)
	if err !=nil{
	return "Error sending Payment"
	}
	wallet.Notes = tmpNotes
	wallet.Notes = append(wallet.Notes, change)
	return "Sent Payment"
	

}

func removeNote(s []note.Note) []note.Note {
    
    return s[:len(s)-1]
}