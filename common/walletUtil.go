package common
import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	//"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/wallet"

)

func SaveWallet(wallet wallet.Wallet, dirname string){

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
		//file = crypto.EncryptWithPublicKey( file, &MyNode.RSAPub )
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}


func ImportWallet(dirname, walletId string) (wallet.Wallet, error){
	
	
	var aWallet wallet.Wallet
    var errWal wallet.Wallet

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
				//theFile:= crypto.DecryptWithPrivateKey(file, &MyNode.RSAPrvt)
				err :=json.Unmarshal([]byte(file), &aWallet)
				
				if err != nil {
					fmt.Println("couldn't unmarshal Wallet", err)
					return errWal, err

				}else{
					return aWallet, err
				}
				
			}else{
				return errWal, err
			}
		}
		
		//fmt.Println( e )
	}else{
		file, err:= ioutil.ReadFile(fileName)
		if err ==nil{
			//theFile:= crypto.DecryptWithPrivateKey(file, &MyNode.RSAPrvt)
			if err !=nil{
				fmt.Println("Could not decrypt file", err)
				return errWal, err
			}else{
				err := json.Unmarshal([]byte(file), &aWallet)
				
				fmt.Println("Unmarshalling Wallet File : ", fileName )
				if err !=nil{
					fmt.Println("Error Unmarshalling Wallet : ", err )
				}else{
					return aWallet, err
				}
			}
		}
	if err != nil {
        fmt.Println("couldn't unmarshal parameters", err)
			return errWal, err
    }
	}


	return errWal, nil
	
}

