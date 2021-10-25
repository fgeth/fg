package node



import (
	"crypto/rsa"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/fgeth/fg/crypto"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/transaction"

)

type Node struct {
	Id				string							//Node Public Key as a string
	Ip				string							//Node IP or fully qualified domain name that can resolve to an IP
	Port			string							//Port that the node is running under
	Path			string							//Path to save data to
	PubKey			*ecdsa.PublicKey				//Nodes Public Key
	PKStr			string							//Node PublicKey as string	
	PrvKey			*ecdsa.PrivateKey				//Nodes Private Key
	PRKStr			string							//Node PrivateKey as string	
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	//NumNodes		uint32							//Tracks Number of Block nodes that have submited Txs
	//Comms			Comm							//Node RSA Keys
	Items			Selling							//Items that are for sell
	
}

type SNode struct {
	Id				string							//Node Public Key as a string
	Ip				string							//Node IP or fully qualified domain name that can resolve to an IP
	Port			string							//Port that the node is running under
	Path			string							//Path to save data to
	PKStr			string							//Node PublicKey as string	
	PRKStr			string							//Node PrivateKey as string	
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	//NumNodes		uint32							//Tracks Number of Block nodes that have submited Txs
	//Comms			Comm							//Node RSA Keys
	Items			Selling							//Items that are for sell
	
}





type Comm struct{

	RsaPrvKeys		map[rsa.PublicKey]rsa.PrivateKey	//index is the RSA publicKey

	
}

type Selling struct {
	Item			map[string]item.Item								 //Index of Item Id and the Item
	Tx				map[string][]transaction.BaseTransaction			//Index is Item Id and array of Debit transactions 
	Keys			map[string][]*ecdsa.PrivateKey						//Index is Item Id and array of private keys for the transaction
	
}
func (node *Node) SNode() SNode{
	var snode SNode
		snode.Id =node.Id
		snode.Ip = node.Ip
		snode.Port = node.Port
		snode.Path = node.Path
		snode.PKStr = node.PKStr
		snode.PRKStr = node.PRKStr
		snode.Writer = node.Writer
		snode.Leader = node.Leader
		return snode
}
func (snode *SNode) Node() Node{
	var node Node
		node.Id =snode.Id
		node.Ip = snode.Ip
		node.Port = snode.Port
		node.Path = snode.Path
		node.PKStr = snode.PKStr
		node.PRKStr = snode.PRKStr
		node.Writer = snode.Writer
		node.Leader = snode.Leader
		return node
}

func (node *Node) SaveNode(dirname string){
	path :=filepath.Join(dirname, "node")
	 //fmt.Println("Path ", path)
	_, err := os.Stat(path)
    if err !=nil {
		fmt.Println("error ", err)
		err := os.Mkdir(dirname, 0755)
		if err !=nil{
			fmt.Println("failed to make root directory", err)
			dirname, _ := os.UserHomeDir()
			node.Path = dirname
			path =filepath.Join(dirname, "node")
			_, err = os.Stat(path)
			if err !=nil{
				err = os.Mkdir(dirname, 0755)
				if err !=nil{
					fmt.Println("failed to make root directory", err)
				}
				err = os.Mkdir(path, 0755)
				if err !=nil{
					fmt.Println("failed to make node directory", err)
				}
			}
		}else{
			err = os.Mkdir(filepath.Join(dirname, "node"), 0755)
			if err !=nil{
				fmt.Println("failed to make node directory", err)
			}
		}

    }
  
	fileName := filepath.Join(path, "node.json")
	snode:=node.SNode()
	file, _ := json.MarshalIndent(snode, "", " ")
	
	//file, _ := json.Marshal(node)
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}

func ImportNode(dirname string) Node{
	
	
	var snode SNode
    
   //fmt.Println( dirname )
	path :=filepath.Join(dirname, "node")

	fileName := filepath.Join(path, "node.json")
	//fmt.Println("File Name : ", fileName )
	_, e := os.Stat(fileName)
	if e != nil{
		dirname, _ := os.UserHomeDir()
		path :=filepath.Join(dirname, "node")
		fileName := filepath.Join(path, "node.json")
		//fmt.Println("File Name : ", fileName )
		_, e1 := os.Stat(fileName)
		
		if e1 != nil{
			fmt.Println( e1 )
		}else{
			file, _ := ioutil.ReadFile(fileName)
			//fmt.Println("Unmarshalling File : ", fileName )
			err :=json.Unmarshal(file, &snode)
			
			if err != nil {
				fmt.Println("couldn't unmarshal parameters", err)
	

			}
		}
		
		//fmt.Println( e )
	}else{
		file, _ := ioutil.ReadFile(fileName)
		
		err := json.Unmarshal(file, &snode)
		//fmt.Println("Unmarshalling File : ", fileName )
	if err != nil {
        fmt.Println("couldn't unmarshal parameters", err)

    }
	}
		node :=snode.Node()
	//fmt.Println("Pub Key Str",  node.PKStr)
	//fmt.Println("Private Key Str",  node.PRKStr)
	//fmt.Println("Node.Id ", node.Id)
	if node.PRKStr !=""{
		node.PrvKey, node.PubKey  = crypto.Decode(node.PRKStr, node.PKStr)
		}
	return node
	
}


