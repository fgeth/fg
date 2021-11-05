package item

import(
	"crypto/rsa"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/transaction"

)

type Items struct {
	Item				map[string]Item					//Index is Item Id
}
type Item struct {
	Id					string
	ProductId			string
	Title				string
	Description			string
	Country				string
	State				string
	City				string
	Image				string
	Amount				float64							//In virtual Dollars
	Qty					uint32
	Color				string
	Weight				Weight
	Height				Size
	Length				Size
	Width				Size
	Tx					TX		//Index is Item Id and array of Debit transactions 
	Seller				rsa.PublicKey
	Buyer				rsa.PublicKey
	Comm				rsa.PrivateKey
	WalletId			string
	Auth				string
	Address				string			//payout address
}



type Weight struct {
	Unit 		string				//oz, lbs, etc..
	Amt			float64

}

type Size struct {
	Unit 		string				//in,ft, mm, meter, etc..
	Amt			float64

}

type TX struct{
Tx	map[string][]transaction.BaseTransaction
}


func CreateItem(id string, productId string, title string, description string, country string, state string, city string, image string, amount float64, qty uint32, color string, weight Weight, height Size,length Size, width Size, tx TX, seller rsa.PublicKey,comm	rsa.PrivateKey, walletId string, auth string, address string) Item{
	
	
	return Item{id,productId,title,description,country,state,city,image,amount,qty,color,weight,height,length,width,tx, seller,seller,comm, walletId, auth, address}

}

func (item *Item) SaveItem(dirname string){
   fmt.Println("Saving Item To ",dirname)
	path :=filepath.Join(dirname, "items")
	 
	_, err := os.Stat(path)
	
    if os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		fmt.Println("Creating Root Directory", err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println("Creating Items Directory", err2)
		
    }
	fileName := filepath.Join(path, item.Id)
	fmt.Println(fileName)
	file, _ := json.MarshalIndent(item, "", " ")
 
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("Error saving Item: ",err)
	}
	if item.Country !=""{
			path =filepath.Join(path , item.Country )
			err2 := os.Mkdir(path, 0755)
			fmt.Println("Creating Country Folder", err2)
			if item.State !=""{
				path =filepath.Join(path , item.State )
				err2 = os.Mkdir(path, 0755)
				fmt.Println("Creating State Folder",err2)
				if item.City !=""{
					path =filepath.Join(path , item.City )
					err2 = os.Mkdir(path, 0755)
					fmt.Println("Creating City Folder",err2)
					fmt.Println(err2)
				
				}
			
			}
		}
		if item.ProductId !=""{
				path =filepath.Join(path , item.ProductId )
				err2 := os.Mkdir(path, 0755)
				fmt.Println("Creating ProductId Folder",err2)
				fmt.Println(err2)
			
			}
	fileName = filepath.Join(path, item.Id)
	fmt.Println(fileName)
	
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("Error saving Item: ",err)
	}

}

func ImportItem(id, dirname string) Item{
	path :=filepath.Join(dirname, "items")
	 _, err := os.Stat(path)
    if err != nil {
        fmt.Println( "error access Item directory", err )
    }
	
	path =filepath.Join(path, id )
	file, _ := ioutil.ReadFile(path)
	var item Item
	_ = json.Unmarshal([]byte(file), &item)
	
	return item
}

func (item Item) Buy() Buy{
	var buyItem Buy
	buyItem.Id = item.Id
	buyItem.Amount = item.Amount
	buyItem.Seller = item.Seller
	buyItem.ProductId = item.ProductId
	buyItem.Country = item.Country
    buyItem.State =item.State
	buyItem.City = item.City
	buyItem.Address = item.Address
	 return buyItem
}



func (item Item) ItemHash() crypto.Hash{

	kh :=crypto.NewKeccakState()
	
	json , _:= json.Marshal(item)
	
	return crypto.HashData(kh, []byte(json))


}