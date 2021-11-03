package common

import(
	"fmt"
	"io/ioutil"
	"path/filepath"
	"github.com/fgeth/fg/item"

)



func SellItem(item item.Item) {
	//prvKey := crypto.GenerateRSAKey()
	//MyNode.Comms.RsaPrvKeys[prvKey.PublicKey] = prvKey
	//item.Seller = prvKey.PublicKey
	//MyNode.Items.Item = append(MyNode.Items.Item, item)

	
}

func AllItemsInDir() {
	dir :=filepath.Join(MyNode.Path, "fg", "items")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		  Items[f.Name()] = item.ImportItem(f.Name(), MyNode.Path)
	}
   
	  
}