package node



import (
	"crypto/rsa"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/transaction"

)

type Node struct {
	Id				string							//Node Public Key as a string
	Ip				string							//Node IP or fully qualified domain name that can resolve to an IP
	Port			string							//Port that the node is running under
	Path			string							//Path to save data to
	PubKey			*ecdsa.PublicKey				//Nodes Public Key
	PrvKey			*ecdsa.PrivateKey				//Nodes Private Key
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	//NumNodes		uint32							//Tracks Number of Block nodes that have submited Txs
	//Comms			Comm							//Node RSA Keys
	//Items			Selling							//Items that are for sell
	
}

type Nodes struct {
	Node 			map[string]Node  				//All active and inactive Nodes.  Easy to get PublicKey using string of public key as map index 

}

type Comm struct{

	RsaPrvKeys		map[rsa.PublicKey]rsa.PrivateKey	//index is the RSA publicKey

	
}

type Selling struct {
	Item			[]item.Item							//Array of Items
	Tx				transaction.Transaction				//Tranaction shell just has the Debit BaseTransactions
	
}

func (node *Node) SaveNode(dirname string){
	path :=filepath.Join(dirname, "fg", "node")
	 fmt.Println("Path ", path)
	_, err := os.Stat(path)
    if err !=nil {
		fmt.Println("error ", err)
		err := os.Mkdir(filepath.Join(dirname, "fg"), 0755)
		if err !=nil{
			fmt.Println("failed to make root directory", err)
		}
        err1 := os.Mkdir(filepath.Join(dirname, "fg", "node"), 0755)
		if err1 !=nil{
			fmt.Println("failed to make node directory", err1)
		}

    }
  
	fileName := filepath.Join(path, "node.json")
	file, _ := json.MarshalIndent(node, "", " ")
	//file, _ := json.Marshal(node)
	_ = ioutil.WriteFile(fileName, file, 0644)

}

func ImportNode(dirname string) Node{
	
	var node Node
    
    fmt.Println( dirname )
	path :=filepath.Join(dirname, "fg", "node")

	fileName := filepath.Join(path, "node.json")
	fmt.Println("File Name : ", fileName )
	myfile, e := os.Stat(fileName)
	if e != nil{
	  fmt.Println( e )
	}else{
		file, _ := ioutil.ReadFile(myfile.Name())
		
		_ = json.Unmarshal([]byte(file), &node)
		
		
	}
	fmt.Println("Node.Id ", node.Id)
	return node
}








