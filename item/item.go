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
	Amount				float64
	Qty					uint32
	Color				string
	Weight				Weight
	Height				Size
	Length				Size
	Width				Size
	Tx					TX		//Index is Item Id and array of Debit transactions 
	Seller				rsa.PublicKey
	Buyer				rsa.PublicKey
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


func CreateItem(id string, productId string, title string, description string, amount float64, qty uint32, color string, weight Weight, height Size,length Size, width Size, tx TX, seller rsa.PublicKey) Item{
	
	
	return Item{id,productId,title,description,amount,qty,color,weight,height,length,width,tx, seller,seller}

}

func (item *Item) SaveItem(dirname string){
   
	path :=filepath.Join(dirname, "fg", "items")
	 
	_, err := os.Stat(path)
	
    if os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(dirname, "fg"), 0755)
		fmt.Println(err)
		err2 := os.Mkdir(path, 0755)
		fmt.Println(err2)
    }
	
	fileName := filepath.Join(path,item.Id)
	fmt.Println(fileName)
	file, _ := json.MarshalIndent(item, "", " ")
 
	_ = ioutil.WriteFile(fileName, file, 0644)

}

func ImportItem(id string) Item{
	dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
	
	path :=filepath.Join(dirname, "fg", "items", id )
	file, _ := ioutil.ReadFile(path)
	var item Item
	_ = json.Unmarshal([]byte(file), &item)
	
	return item
}

func (item Item) ItemHash() crypto.Hash{

	kh :=crypto.NewKeccakState()
	
	json , _:= json.Marshal(item)
	
	return crypto.HashData(kh, []byte(json))


}