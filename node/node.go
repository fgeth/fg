package node



import (
	"crypto/rsa"
	"crypto/ecdsa"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/transaction"

)

type Node struct {
	Id				string							//Node Public Key as a string
	Ip				string							//Node IP or fully qualified domain name that can resolve to an IP
	Port			string							//Port that the node is running under
	PubKey			*ecdsa.PublicKey				//Nodes Public Key
	PrvKey			*ecdsa.PrivateKey				//Nodes Private Key
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	NumNodes		uint32							//Tracks Number of Block nodes that have submited Txs
	Comms			Comm							//Node RSA Keys
	Items			Selling							//Items that are for sell
	
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

func (node *Node) SaveNode(){

}









