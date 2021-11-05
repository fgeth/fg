package node

import(
	"encoding/json"
    "fmt"
	"io/ioutil"
	"os"
	"path/filepath"

)


type SNode struct {
	Id				uint64							//Node Ring Location
	Ip				string							//Node Onion Address
	Port			string							//Port that the node is running under
	Path			string							//Path to save data to
	PKStr			string							//Node PublicKey as string	
	PRKStr			string							//Node PrivateKey as string	
	Address			string							//File Name for Password Protected Stored Private Key	
	WalletId		string							//Address for wallet
}



//From Saved File to Node
func (snode *SNode) Node() Node{
	var node Node
		node.Id=snode.Id
		node.Ip = snode.Ip
		node.Port = snode.Port
		node.Path = snode.Path
		node.PKStr = snode.PKStr
		node.PRKStr = snode.PRKStr
		node.Address = snode.Address
		node.WalletId = snode.WalletId
		return node
}



func (node *SNode) SaveNodeOne(dirname string){
	path :=filepath.Join(dirname, "node")
	 //fmt.Println("Path ", path)
	_, err := os.Stat(path)
    
	fileName := filepath.Join(path, "nodeOne.json")
	
	file, _ := json.MarshalIndent(node, "", " ")
	
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}

func (node *SNode) SaveNodeTwo(dirname string){
	path :=filepath.Join(dirname, "node")
	 //fmt.Println("Path ", path)
	_, err := os.Stat(path)
      
	fileName := filepath.Join(path, "nodeTwo.json")

	file, _ := json.MarshalIndent(node, "", " ")

	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}