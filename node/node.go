package node



import (
	//"crypto/rsa"
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
	PRKStr			string							//Node PrivateKey as string	
	Address			string							//File Name for Password Protected Stored Private Key
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	Keys			Key								//Array of Private Keys as Addresses
	Comms			Comm							//Node RSA Keys
	WalletId		string							//Address
	
	
}


type Key struct {

	Key 		[]string	//Array of Addresses which is the file location of the Private Keys
}

type Comm struct{

	Rsa			[]string	//Array of Addresses which is the file location of the Private Keys

	
}
//Node that can be saved to file
func (node *Node) SNode() SNode{
	var snode SNode
		snode.Id =node.Id
		snode.Ip = node.Ip
		snode.Port = node.Port
		snode.Path = node.Path
		snode.PKStr = node.PKStr
		snode.PRKStr = node.PRKStr
		snode.Address = node.Address
		snode.WalletId = node.WalletId
		return snode
}

//Public Node
func (node *Node) PNode() PNode{
	var pnode PNode
		pnode.Id =node.Id
		pnode.Ip = node.Ip
		pnode.Port = node.Port
		pnode.PKStr = node.PKStr

		return pnode
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
	file, err := json.MarshalIndent(snode, "", " ")
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

