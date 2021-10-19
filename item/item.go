package item

import(
	"math/big"
	"crypto/rsa"
	"path/filepath"

)

type Items struct {
	Item				map[string]Item					//Index is Item Id
}
type Item struct {
	Id					string
	ProductId			string
	Title				string
	Description			string
	Amount				*big.Int
	Qty					uint32
	Color				string
	Weight				Weight
	Height				Size
	Length				Size
	Width				Size
	Seller				rsa.PublicKey
	Buyer				rsa.PublicKey
}

type Weight struct {
	Unit 		string				//oz, lbs, etc..
	Amount		float32

}

type Size struct {
	Unit 		string				//in,ft, mm, meter, etc..
	Amount		float32

}

func CreateItem(id string, productId string, title string, description string, amount *big.Int, qty uint32, color string, weight weight, height size,length size, width size, seller rsa.PublicKey) Item{
	
	
	return Item{id,productId,title,description,amount,qty,color,weight,height,length,width,seller,seller}

}

func (item *Item) SaveItem(){
    dirname, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }
 
	path :=filepath.Join(dirname, "fg", "items")
	 
	folderInfo, err := os.Stat(path)
	if folderInfo.Name() !="" {
			fmt.Println("")
	}
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
					  Items = Items{map[uint64]Item{info.Name(), item.ImportItem(info.Name)}}
              }
             
      })
	  
}