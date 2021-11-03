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


)

type Node struct {
	Id				uint64							//Node Ring Location
	Ip				string							//Node Onion Address
	Port			string							//Port that the node is running under
	Path			string							//Path to save data to
	PubKey			*ecdsa.PublicKey				//Nodes Public Key
	PKStr			string							//Node PublicKey as string	
	PrvKey			*ecdsa.PrivateKey				//Nodes Private Key
	Address			string							//File Name for Password Protected Stored Private Key
	PRKStr			string							//Node PrivateKey as string	
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	Addresses		[]string						//Array of Public Keys as Addresses
	PubKeys			[]*ecdsa.PublicKey				//Array of Public Keys
	PrvtKeys		[]*ecdsa.PrivateKey				//Array of Private Keys
	Comms			Comm							//Node RSA Keys
	
	
}

type SNode struct {
	Id				uint64							//Node Ring Location
	Ip				string							//Node Onion Address
	Port			string							//Port that the node is running under
	PKStr			string							//Node PublicKey as string	
	Address			string							//File Name for Password Protected Stored Private Key	
	
}

type RNode struct {
	Id				uint64							//Node Public Key as a uint64
	PKStr			string							//Node PublicKey as string	
	
	
}



type Comm struct{

	RsaPrvKeys		map[rsa.PublicKey]rsa.PrivateKey	//index is the RSA publicKey

	
}

func (node *Node) SNode() SNode{
	var snode SNode
		snode.Id =node.Id
		snode.Ip = node.Ip
		snode.Port = node.Port
		snode.PKStr = node.PKStr
		return snode
}
func (snode *SNode) Node() Node{
	var node Node
		node.Id=snode.Id
		node.Ip = snode.Ip
		node.Port = snode.Port
		node.PKStr = snode.PKStr
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
	//snode:=node.SNode()
	file, err := json.MarshalIndent(node, "", " ")
	if err !=nil{
		fmt.Println("Error Marshalling Node :", err)
	}
	fmt.Println("The Node Marshalled: ", file)
	//file, _ := json.Marshal(node)
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
	

}

func ImportNode(dirname string) (Node, error){
	
	
	var snode SNode
    var errNode Node
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
			return errNode, e1
			
		}else{
			file, _ := ioutil.ReadFile(fileName)
			//fmt.Println("Unmarshalling File : ", fileName )
			err :=json.Unmarshal(file, &snode)
			
			if err != nil {
				fmt.Println("couldn't unmarshal parameters", err)
				return errNode, err

			}
		}
		
		//fmt.Println( e )
	}else{
		file, _ := ioutil.ReadFile(fileName)
		
		err := json.Unmarshal(file, &snode)
		//fmt.Println("Unmarshalling File : ", fileName )
	if err != nil {
        fmt.Println("couldn't unmarshal parameters", err)
			return errNode, err
    }
	}
		node :=snode.Node()
	//fmt.Println("Pub Key Str",  node.PKStr)
	//fmt.Println("Private Key Str",  node.PRKStr)
	//fmt.Println("Node.Id ", node.Id)
	if node.PRKStr !=""{
		node.PrvKey, node.PubKey  = crypto.Decode(node.PRKStr, node.PKStr)
		}
	return node, nil
	
}

func (node *Node) GetNodes(){

}

func (node *Node) RegisterNode(PubKey string){

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
