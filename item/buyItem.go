package item

import(
	"crypto/rsa"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)



type Buy struct {
	Id					string
	Amount				float64
	Seller				rsa.PublicKey
	Buyer				rsa.PublicKey
	ProductId			string
	Country				string
	State				string
	City				string
	Address				string			//payout address
	WalletId			string
	Password			string
}






func (buyItem *Buy) ImportItem( dirname string) Item{
	path :=filepath.Join(dirname, "items")
	 _, err := os.Stat(path)
    if err != nil {
        fmt.Println( "error access Item directory", err )
    }
	if buyItem.Country !=""{
		path =filepath.Join(path, buyItem.Country )
		if buyItem.State !=""{
			path =filepath.Join(path, buyItem.State )
			if buyItem.City !=""{
				path =filepath.Join(path, buyItem.City )
			}
		}
	}
	if buyItem.ProductId !=""{
		path =filepath.Join(path, buyItem.ProductId )
	}
	path =filepath.Join(path, buyItem.Id )
	file, _ := ioutil.ReadFile(path)
	var item Item
	_ = json.Unmarshal([]byte(file), &item)
	
	return item
}